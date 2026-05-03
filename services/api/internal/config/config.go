package config

import "os"

type Config struct {
	Env               string
	Host              string
	Port              string
	DatabaseURL       string
	RedisURL          string
	RabbitMQURL       string
	S3Endpoint        string
	S3Bucket          string
	GoogleClientID    string
	FacebookLoginFlag string
}

func Load() Config {
	return Config{
		Env:               get("APP_ENV", "development"),
		Host:              get("API_HOST", "0.0.0.0"),
		Port:              get("API_PORT", "8080"),
		DatabaseURL:       get("DATABASE_URL", ""),
		RedisURL:          get("REDIS_URL", ""),
		RabbitMQURL:       get("RABBITMQ_URL", ""),
		S3Endpoint:        get("S3_ENDPOINT", ""),
		S3Bucket:          get("S3_BUCKET", ""),
		GoogleClientID:    get("GOOGLE_OAUTH_CLIENT_ID", ""),
		FacebookLoginFlag: get("FACEBOOK_LOGIN_ENABLED", "false"),
	}
}

func get(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
