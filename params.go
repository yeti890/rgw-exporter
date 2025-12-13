package main

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AccessKey string
	SecretKey string

	// Internal RGW Admin endpoint
	Endpoint string

	Region string

	ClusterName string

	// Public S3 endpoint (used as label "endpoint")
	PubEndpoint string

	ListenIP   string
	ListenPort int

	UsageCollectorInterval   int
	BucketsCollectorInterval int
	UsersCollectorInterval   int

	RGWConnectionTimeout int
	StartDelay           int
	Insecure             bool
	SkipWithoutBucket    bool

	UsersCollectorEnable bool
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func loadConfig() (*Config, error) {
	cfg := &Config{
		AccessKey: getEnv("ACCESS_KEY", ""),
		SecretKey: getEnv("SECRET_KEY", ""),

		Endpoint: getEnv("RGW_ENDPOINT", ""),

		Region:      getEnv("REGION", ""),
		ClusterName: getEnv("CLUSTER_NAME", ""),

		PubEndpoint: getEnv("PUB_ENDPOINT", ""),

		ListenIP:   getEnv("LISTEN_IP", "127.0.0.1"),
		ListenPort: getEnvInt("LISTEN_PORT", 9240),

		UsageCollectorInterval:   getEnvInt("USAGE_COLLECTOR_INTERVAL", 30),
		BucketsCollectorInterval: getEnvInt("BUCKETS_COLLECTOR_INTERVAL", 300),
		UsersCollectorInterval:   getEnvInt("USERS_COLLECTOR_INTERVAL", 600),

		RGWConnectionTimeout: getEnvInt("RGW_CONNECTION_TIMEOUT", 600),
		StartDelay:           getEnvInt("START_DELAY", 30),

		Insecure:          getEnvBool("INSECURE", false),
		SkipWithoutBucket: getEnvBool("SKIP_WITHOUT_BUCKET", false),

		UsersCollectorEnable: getEnvBool("USERS_COLLECTOR_ENABLE", false),
	}

	// ---- Required fields validation ----
	if cfg.AccessKey == "" {
		return nil, fmt.Errorf("ACCESS_KEY is required")
	}
	if cfg.SecretKey == "" {
		return nil, fmt.Errorf("SECRET_KEY is required")
	}
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("RGW_ENDPOINT is required")
	}

	// PUB_ENDPOINT technically can be empty, but we strongly recommend setting it
	// to make label "endpoint" meaningful. We keep it non-fatal to avoid breaking
	// minimal lab setups.
	return cfg, nil
}
