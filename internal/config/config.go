package config

import (
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT"     envDefault:"5432"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME,required"`
	SSLMode  string `env:"DB_SSLMODE"  envDefault:"disable"`
}

type Config struct {
	Port           int `env:"PORT"           envDefault:"8080"`
	DB             Database
	GoogleClientID string `env:"GOOGLE_CLIENT_ID,required"`
	JWTSecret      string `env:"JWT_SECRET"`
	AllowedOrigins string `env:"ALLOWED_ORIGINS,required"`
	AppEnv         string `env:"APP_ENV"         envDefault:"development"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) DSN() string {
	u := &url.URL{
		Scheme: "postgres",
		Host:   net.JoinHostPort(c.DB.Host, strconv.Itoa(c.DB.Port)),
		Path:   c.DB.Name,
	}
	if c.DB.Password != "" {
		u.User = url.UserPassword(c.DB.User, c.DB.Password)
	} else {
		u.User = url.User(c.DB.User)
	}

	q := u.Query()
	q.Set("sslmode", c.DB.SSLMode)
	u.RawQuery = q.Encode()

	return u.String()
}

func (c *Config) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("config: PORT out of range: %d", c.Port)
	}

	if c.DB.Host == "" {
		return fmt.Errorf("config: DB_HOST is required")
	}
	if c.DB.Port < 1 || c.DB.Port > 65535 {
		return fmt.Errorf("config: DB_PORT out of range: %d", c.DB.Port)
	}
	if c.DB.User == "" {
		return fmt.Errorf("config: DB_USER is required")
	}
	if c.DB.Name == "" {
		return fmt.Errorf("config: DB_NAME is required")
	}
	if c.DB.SSLMode == "" {
		return fmt.Errorf("config: DB_SSLMODE is required")
	}

	if _, err := pgxpool.ParseConfig(c.DSN()); err != nil {
		return fmt.Errorf("config: database connection is not valid: %w", err)
	}

	if c.GoogleClientID == "" {
		return fmt.Errorf("config: GOOGLE_CLIENT_ID is required")
	}
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("config: JWT_SECRET must be at least 32 chars")
	}
	if c.AllowedOrigins == "" {
		return fmt.Errorf("config: ALLOWED_ORIGINS is required")
	}

	return nil
}
