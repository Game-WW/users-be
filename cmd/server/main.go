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

	srv, err := buildServer(cfg, pool)
	if err != nil {
		log.Fatalf("server: %v", err)
	}

	grpcServer := grpc.NewServer()
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
