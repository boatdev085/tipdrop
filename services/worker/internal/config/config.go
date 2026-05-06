package config

import "os"

type Config struct {
	Env         string
	DatabaseURL string
	RedisURL    string
	RabbitMQURL string
}

func Load() Config {
	return Config{
		Env:         get("APP_ENV", "development"),
		DatabaseURL: get("DATABASE_URL", ""),
		RedisURL:    get("REDIS_URL", ""),
		RabbitMQURL: get("RABBITMQ_URL", ""),
	}
}

func get(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
