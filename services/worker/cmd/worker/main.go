package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/boatdev085/tipdrop/services/worker/internal/config"
)

func main() {
	cfg := config.Load()
	log.Printf("tipdrop worker starting env=%s rabbitmq_configured=%t", cfg.Env, cfg.RabbitMQURL != "")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = shutdownCtx

	log.Println("tipdrop worker stopped")
}
