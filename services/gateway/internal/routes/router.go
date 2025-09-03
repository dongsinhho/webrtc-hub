package routes

import (
	"net/http"

	"github.com/dongsinhho/webrtc-hub/services/gateway/internal/config"
	"github.com/dongsinhho/webrtc-hub/services/gateway/internal/middleware"
	"github.com/dongsinhho/webrtc-hub/services/gateway/internal/proxy"
	"github.com/dongsinhho/webrtc-hub/services/gateway/internal/telemetry"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func New(cfg config.Config) (http.Handler, error) {
	r := chi.NewRouter()

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // or parse cfg.AllowedOrigins
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	// Recovery & logging & rate limit
	r.Use(middleware.Recover)
	r.Use(middleware.Logging)
	r.Use(httprate.Limit(200, telemetry.Window()))

	// Health
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	// Metrics
	r.Method("GET", "/metrics", telemetry.MetricsHandler())

	// JWT for protected routes
	protected := chi.NewRouter()
	protected.Use(middleware.JWTAuth(middleware.AuthConfig{JWKSURL: cfg.JWTPublicJWKS}))

	// Reverse proxy: /auth/* → auth-svc
	if rp, err := proxy.NewReverseProxy(cfg.AuthBaseURL, "/auth"); err == nil {
		protected.Mount("/auth", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { rp.ServeHTTP(w, r) }))
	}

	// Reverse proxy: /rooms/* → room-svc
	if rp, err := proxy.NewReverseProxy(cfg.RoomBaseURL, "/rooms"); err == nil {
		protected.Mount("/rooms", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { rp.ServeHTTP(w, r) }))
	}

	// WebSocket signaling passthrough: /ws/signaling → signaling-svc
	if rp, err := proxy.NewReverseProxy(cfg.SignalingURL, ""); err == nil {
		// assuming backend has /ws/signaling path too; otherwise mount raw hijack
		r.Mount("/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { rp.ServeHTTP(w, r) }))
	}

	// Mount protected API subtree
	r.Mount("/api", protected)

	return r, nil
}
