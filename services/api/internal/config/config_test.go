package config

import "testing"

func TestMissingRequired(t *testing.T) {
	cfg := Config{
		DatabaseURL:       "postgres://local",
		RedisURL:          "redis://local",
		RabbitMQURL:       "amqp://local",
		S3Endpoint:        "http://localhost:9000",
		S3Region:          "us-east-1",
		S3Bucket:          "tipdrop-local",
		S3AccessKeyID:     "tipdrop",
		S3SecretAccessKey: "tipdrop-secret",
		JWTAccessSecret:   "access",
		JWTRefreshSecret:  "refresh",
	}

	if missing := cfg.MissingRequired(); len(missing) != 0 {
		t.Fatalf("expected no missing config, got %v", missing)
	}
}

func TestMissingRequiredReturnsStableOrder(t *testing.T) {
	missing := (Config{}).MissingRequired()
	want := []string{
		"DATABASE_URL",
		"REDIS_URL",
		"RABBITMQ_URL",
		"S3_ENDPOINT",
		"S3_REGION",
		"S3_BUCKET",
		"S3_ACCESS_KEY_ID",
		"S3_SECRET_ACCESS_KEY",
		"JWT_ACCESS_SECRET",
		"JWT_REFRESH_SECRET",
	}

	if len(missing) != len(want) {
		t.Fatalf("expected %d missing keys, got %d: %v", len(want), len(missing), missing)
	}
	for index := range want {
		if missing[index] != want[index] {
			t.Fatalf("missing[%d] = %q, want %q", index, missing[index], want[index])
		}
	}
}
