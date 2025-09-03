package config

import (
	"log"
	"os"
)

type Config struct {
	HTTPAddr       string
	AuthBaseURL    string // http://auth-svc:8082
	RoomBaseURL    string // http://room-svc:8081
	SignalingURL   string // http://signaling-svc:8080 (WS capable)
	JWTPublicJWKS  string // URL or path to JWKS (or PEM)
	AllowedOrigins string // CORS, comma separated
}

func mustEnv(key, val string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	if val != "" {
		return val
	}
	log.Fatalf("missing required env %s", key)
	return ""
}

func Load() Config {
	return Config{
		HTTPAddr:       mustEnv("GATEWAY_HTTP_ADDR", ":8080"),
		AuthBaseURL:    mustEnv("AUTH_BASE_URL", "http://auth-svc:8082"),
		RoomBaseURL:    mustEnv("ROOM_BASE_URL", "http://room-svc:8081"),
		SignalingURL:   mustEnv("SIGNALING_URL", "http://signaling-svc:8080"),
		JWTPublicJWKS:  os.Getenv("JWT_JWKS_URL"),
		AllowedOrigins: os.Getenv("CORS_ALLOWED_ORIGINS"),
	}
}
