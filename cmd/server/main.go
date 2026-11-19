package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	userspb "github.com/Game-WW/gateway-be/gen/users"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/Game-WW/users-be/internal/infrastructure/database"
	"github.com/Game-WW/users-be/internal/infrastructure/s3"
)

const dbConnectTimeout = 10 * time.Second

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("config: .env not found, using environment")
	}

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	connectCtx, cancel := context.WithTimeout(ctx, dbConnectTimeout)
	pool, err := database.NewPool(connectCtx, cfg.DSN)
	cancel()
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	s3Client, err := s3.NewClient(
		cfg.S3Endpoint,
		cfg.S3AccessKey,
		cfg.S3SecretKey,
		cfg.S3Region,
	)
	if err != nil {
		log.Fatalf("s3: %v", err)
	}

	s3Ctx, cancelS3 := context.WithTimeout(ctx, 5*time.Second)
	err = s3.EnsureBucket(s3Ctx, s3Client, cfg.S3Bucket, cfg.S3Region)
	cancelS3()
	if err != nil {
		log.Fatalf("s3: ensure bucket %q: %v", cfg.S3Bucket, err)
	}

	storage := s3.NewStorage(s3Client, cfg.S3Bucket)
	srv, err := buildServer(cfg, pool, storage)
	if err != nil {
		log.Fatalf("server: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.MaxRecvMsgSize(6 << 20))
	userspb.RegisterUsersServiceServer(grpcServer, srv)

	addr := ":" + strconv.Itoa(cfg.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("server: listen %s: %v", addr, err)
	}

	go func() {
		<-ctx.Done()
		log.Println("server: shutting down")
		grpcServer.GracefulStop()
	}()

	log.Printf("server: grpc listening on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("server: serve: %v", err)
	}
}
