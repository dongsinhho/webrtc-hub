package main

import (
	"log"
	"net/http"

	"github.com/dongsinhho/webrtc-hub/services/gateway/internal/config"
	"github.com/dongsinhho/webrtc-hub/services/gateway/internal/routes"
)

func main() {
	cfg := config.Load()

	h, err := routes.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("gateway listening on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, h); err != nil {
		log.Fatal(err)
	}
}
