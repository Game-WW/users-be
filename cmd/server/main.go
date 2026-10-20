package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/Game-WW/users-be/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("config: .env not found, using defaults")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("fatal: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		log.Fatalf("database: invalid configuration: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		log.Fatalf("database: unable to connect to pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("database: ping failed: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName:      "users-be",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowHeaders:     "Authorization, Content-Type",
		AllowCredentials: true,
	}))

	port := cfg.Port
	addr := ":" + strconv.Itoa(port)
	if err = app.Listen(addr); err != nil {
		log.Fatalf("server: failed to listen on %s: %v", addr, err)
	}
}
