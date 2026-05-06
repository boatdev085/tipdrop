package httpserver

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boatdev085/tipdrop/services/api/internal/config"
	"github.com/gofiber/fiber/v2"
)

func Run(cfg config.Config) error {
	app := fiber.New(fiber.Config{
		AppName: "tipdrop-api",
	})

	registerRoutes(app, cfg)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Listen(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-stopCh:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(ctx); err != nil {
			return err
		}
		return nil
	}
}

func registerRoutes(app *fiber.App, cfg config.Config) {
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"ok":      true,
			"service": "api",
			"env":     cfg.Env,
		})
	})

	app.Get("/readyz", func(c *fiber.Ctx) error {
		missing := cfg.MissingRequired()
		if len(missing) > 0 {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"ok":      false,
				"missing": missing,
			})
		}

		return c.JSON(fiber.Map{
			"ok": true,
		})
	})

	api := app.Group("/api")
	api.Get("/me", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error": "not_implemented",
		})
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": errors.New("not found").Error(),
		})
	})
}
