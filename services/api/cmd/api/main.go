package main

import (
	"log"

	"github.com/boatdev085/tipdrop/services/api/internal/config"
	"github.com/boatdev085/tipdrop/services/api/internal/httpserver"
)

func main() {
	cfg := config.Load()

	if err := httpserver.Run(cfg); err != nil {
		log.Fatalf("api stopped: %v", err)
	}
}
