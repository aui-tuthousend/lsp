package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv      string
	HTTPAddr    string
	DisableAuth bool

	PostgresDSN string
	RedisAddr   string
	RedisPass   string
	RedisDB     int

	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2Bucket          string
	R2PublicBaseURL   string
	R2Region          string

	MailerDriver  string
	MailgunDomain string
	MailgunAPIKey string
	MailgunSender string

	JWTSecret string
	JWTExpiry time.Duration
}

func Load() (Config, error) {
	redisDB, err := parseInt(env("REDIS_DB", "0"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid REDIS_DB: %w", err)
	}

	disableAuth, err := parseBool(env("DISABLE_AUTH", "false"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid DISABLE_AUTH: %w", err)
	}

	jwtExpiry, err := time.ParseDuration(env("JWT_EXPIRY", "24h"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid JWT_EXPIRY: %w", err)
	}

	cfg := Config{
		AppEnv:            env("APP_ENV", "development"),
		HTTPAddr:          env("HTTP_ADDR", ":8080"),
		DisableAuth:       disableAuth,
		PostgresDSN:       os.Getenv("POSTGRES_DSN"),
		RedisAddr:         env("REDIS_ADDR", "localhost:6379"),
		RedisPass:         os.Getenv("REDIS_PASSWORD"),
		RedisDB:           redisDB,
		R2AccountID:       os.Getenv("R2_ACCOUNT_ID"),
		R2AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		R2SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
		R2Bucket:          os.Getenv("R2_BUCKET"),
		R2PublicBaseURL:   strings.TrimRight(os.Getenv("R2_PUBLIC_BASE_URL"), "/"),
		R2Region:          env("R2_REGION", "auto"),
		MailerDriver:      env("MAILER_DRIVER", "log"),
		MailgunDomain:     os.Getenv("MAILGUN_DOMAIN"),
		MailgunAPIKey:     os.Getenv("MAILGUN_API_KEY"),
		MailgunSender:     os.Getenv("MAILGUN_SENDER"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		JWTExpiry:         jwtExpiry,
	}
	return cfg, cfg.Validate()
}

func (c Config) Validate() error {
	if c.HTTPAddr == "" {
		return errors.New("HTTP_ADDR is required")
	}
	if c.PostgresDSN == "" {
		return errors.New("POSTGRES_DSN is required")
	}
	if c.RedisAddr == "" {
		return errors.New("REDIS_ADDR is required")
	}
	if c.JWTExpiry <= 0 {
		return errors.New("JWT_EXPIRY must be positive")
	}
	if c.MailerDriver != "log" && c.MailerDriver != "mailgun" {
		return errors.New("MAILER_DRIVER must be log or mailgun")
	}
	if !c.DisableAuth && c.JWTSecret == "" {
		return errors.New("JWT_SECRET is required when auth is enabled")
	}
	if c.AppEnv == "production" {
		if c.DisableAuth {
			return errors.New("DISABLE_AUTH cannot be true in production")
		}
		if len(c.JWTSecret) < 32 {
			return errors.New("JWT_SECRET must be at least 32 characters in production")
		}
		if c.MailerDriver == "log" {
			return errors.New("MAILER_DRIVER=log is not allowed in production")
		}
	}
	if c.MailerDriver == "mailgun" {
		if c.MailgunDomain == "" || c.MailgunAPIKey == "" || c.MailgunSender == "" {
			return errors.New("MAILGUN_DOMAIN, MAILGUN_API_KEY, and MAILGUN_SENDER are required")
		}
	}
	if c.HasR2Config() {
		if err := c.ValidateR2(); err != nil {
			return err
		}
	}
	return nil
}

func (c Config) HasR2Config() bool {
	return c.R2AccountID != "" || c.R2AccessKeyID != "" || c.R2SecretAccessKey != "" || c.R2Bucket != "" || c.R2PublicBaseURL != ""
}

func (c Config) ValidateR2() error {
	if c.R2AccountID == "" || c.R2AccessKeyID == "" || c.R2SecretAccessKey == "" || c.R2Bucket == "" || c.R2PublicBaseURL == "" {
		return errors.New("R2_ACCOUNT_ID, R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY, R2_BUCKET, and R2_PUBLIC_BASE_URL are required")
	}
	if c.R2Region == "" {
		return errors.New("R2_REGION is required")
	}
	return nil
}

func parseBool(value string) (bool, error) {
	return strconv.ParseBool(value)
}

func parseInt(value string) (int, error) {
	return strconv.Atoi(value)
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
