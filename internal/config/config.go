package config

import (
	"fmt"
	"os"
	"time"
)

// Config agrupa toda la configuración cargada del entorno.
// JWT_SECRET es obligatorio. Las URLs tienen defaults de dev.
type Config struct {
	// Infra
	DatabaseURL string
	RedisURL    string
	Port        string
	JWTSecret   string

	// Tunables del sistema reactivo
	DashboardTTL         time.Duration // TTL de cache del dashboard (5 min)
	ResolverInterval     time.Duration // ticker del resolver UNKNOWN (30 s)
	WorkerReaperInterval time.Duration // ticker del reaper outbox (60 s)
	SSEHeartbeat         time.Duration // heartbeat SSE (15 s)
}

// Load lee y valida la configuración. Falla si falta JWT_SECRET.
func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL: envOr("DATABASE_URL", "postgres://movank:movank@localhost:5432/movank?sslmode=disable"),
		RedisURL:    envOr("REDIS_URL", "redis://localhost:6379/0"),
		Port:        envOr("PORT", "8080"),
		JWTSecret:   os.Getenv("JWT_SECRET"),

		DashboardTTL:         5 * time.Minute,
		ResolverInterval:     30 * time.Second,
		WorkerReaperInterval: 60 * time.Second,
		SSEHeartbeat:         15 * time.Second,
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET no configurado (obligatorio)")
	}

	return cfg, nil
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
