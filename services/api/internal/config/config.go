package config

import "os"

type Config struct {
	Env                  string
	Host                 string
	Port                 string
	DatabaseURL          string
	RedisURL             string
	RabbitMQURL          string
	S3Endpoint           string
	S3Region             string
	S3Bucket             string
	S3AccessKeyID        string
	S3SecretAccessKey    string
	S3ForcePathStyle     string
	JWTAccessSecret      string
	JWTRefreshSecret     string
	GoogleClientID       string
	FacebookAppID        string
	FacebookAppSecret    string
	FacebookLoginFlag    string
	NextPublicAPIBaseURL string
}

func Load() Config {
	return Config{
		Env:                  get("APP_ENV", "development"),
		Host:                 get("API_HOST", "0.0.0.0"),
		Port:                 get("API_PORT", "8080"),
		DatabaseURL:          get("DATABASE_URL", ""),
		RedisURL:             get("REDIS_URL", ""),
		RabbitMQURL:          get("RABBITMQ_URL", ""),
		S3Endpoint:           get("S3_ENDPOINT", ""),
		S3Region:             get("S3_REGION", ""),
		S3Bucket:             get("S3_BUCKET", ""),
		S3AccessKeyID:        get("S3_ACCESS_KEY_ID", ""),
		S3SecretAccessKey:    get("S3_SECRET_ACCESS_KEY", ""),
		S3ForcePathStyle:     get("S3_FORCE_PATH_STYLE", "false"),
		JWTAccessSecret:      get("JWT_ACCESS_SECRET", ""),
		JWTRefreshSecret:     get("JWT_REFRESH_SECRET", ""),
		GoogleClientID:       get("GOOGLE_OAUTH_CLIENT_ID", ""),
		FacebookAppID:        get("FACEBOOK_OAUTH_APP_ID", ""),
		FacebookAppSecret:    get("FACEBOOK_OAUTH_APP_SECRET", ""),
		FacebookLoginFlag:    get("FACEBOOK_LOGIN_ENABLED", "false"),
		NextPublicAPIBaseURL: get("NEXT_PUBLIC_API_BASE_URL", ""),
	}
}

func (cfg Config) MissingRequired() []string {
	checks := []struct {
		key   string
		value string
	}{
		{key: "DATABASE_URL", value: cfg.DatabaseURL},
		{key: "REDIS_URL", value: cfg.RedisURL},
		{key: "RABBITMQ_URL", value: cfg.RabbitMQURL},
		{key: "S3_ENDPOINT", value: cfg.S3Endpoint},
		{key: "S3_REGION", value: cfg.S3Region},
		{key: "S3_BUCKET", value: cfg.S3Bucket},
		{key: "S3_ACCESS_KEY_ID", value: cfg.S3AccessKeyID},
		{key: "S3_SECRET_ACCESS_KEY", value: cfg.S3SecretAccessKey},
		{key: "JWT_ACCESS_SECRET", value: cfg.JWTAccessSecret},
		{key: "JWT_REFRESH_SECRET", value: cfg.JWTRefreshSecret},
	}

	missing := make([]string, 0)
	for _, check := range checks {
		if check.value == "" {
			missing = append(missing, check.key)
		}
	}
	return missing
}

func get(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
