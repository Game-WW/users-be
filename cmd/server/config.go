package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

const minJWTSecretLen = 32

type config struct {
	GRPCPort           int
	DSN                string
	GoogleClientID     string
	JWTSecret          string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
	GitHubClientID     string
	GitHubClientSecret string
	YandexClientID     string
	YandexClientSecret string
}

func loadConfig() (config, error) {
	var cfg config
	var err error

	if cfg.GRPCPort, err = strconv.Atoi(env("GRPC_PORT", "50051")); err != nil {
		return config{}, fmt.Errorf("GRPC_PORT: %w", err)
	}

	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(env("DB_USER", "postgres"), os.Getenv("DB_PASSWORD")),
		Host:     env("DB_HOST", "localhost") + ":" + env("DB_PORT", "5433"),
		Path:     env("DB_NAME", "users"),
		RawQuery: "sslmode=" + env("DB_SSLMODE", "disable"),
	}
	cfg.DSN = dsn.String()

	cfg.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")

	cfg.GitHubClientID = os.Getenv("GITHUB_CLIENT_ID")
	cfg.GitHubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")

	cfg.YandexClientID = os.Getenv("YANDEX_CLIENT_ID")
	cfg.YandexClientSecret = os.Getenv("YANDEX_CLIENT_SECRET")

	// Должен совпадать с JWT_SECRET в gateway-be.
	if cfg.JWTSecret = os.Getenv("JWT_SECRET"); len(cfg.JWTSecret) < minJWTSecretLen {
		return config{}, fmt.Errorf("JWT_SECRET must be at least %d characters", minJWTSecretLen)
	}

	if cfg.AccessTokenTTL, err = ttl("ACCESS_TOKEN_TTL", "15m"); err != nil {
		return config{}, err
	}
	// Должен совпадать с refreshTokenTTL в gateway-be (internal/proxy/cookies.go).
	if cfg.RefreshTokenTTL, err = ttl("REFRESH_TOKEN_TTL", "720h"); err != nil {
		return config{}, err
	}

	return cfg, nil
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func ttl(key, def string) (time.Duration, error) {
	d, err := time.ParseDuration(env(key, def))
	if err != nil {
		return 0, fmt.Errorf("%s: %w", key, err)
	}
	if d < time.Second {
		return 0, fmt.Errorf("%s must be at least 1s", key)
	}
	return d, nil
}
