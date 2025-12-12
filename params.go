package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Endpoint  string `yaml:"endpoint"`
	Region    string `yaml:"region"`

	ClusterName string `yaml:"cluster_name"`

	PubEndpoint string `yaml:"pub_endpoint"`

	ListenIP   string `yaml:"listen_ip"`
	ListenPort int    `yaml:"listen_port"`

	UsageCollectorInterval   int  `yaml:"usage_collector_interval"`
	BucketsCollectorInterval int  `yaml:"buckets_collector_interval"`
	UsersCollectorInterval   int  `yaml:"users_collector_interval"`
	RGWConnectionTimeout     int  `yaml:"rgw_connection_timeout"`
	StartDelay               int  `yaml:"start_delay"`
	Insecure                 bool `yaml:"insecure"`
	SkipWithoutBucket        bool `yaml:"skip_without_bucket"`

	UsersCollectorEnable bool `yaml:"users_collector_enable"`
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
	config := &Config{
		AccessKey: getEnv("ACCESS_KEY", ""),
		SecretKey: getEnv("SECRET_KEY", ""),

		// Public RGW endpoint (Config.Endpoint)
		Endpoint: getEnv("RGW_ENDPOINT", ""),

		Region: getEnv("REGION", ""),

		ClusterName: getEnv("CLUSTER_NAME", ""),

		PubEndpoint: getEnv("PUB_ENDPOINT", ""),

		ListenIP:   getEnv("LISTEN_IP", "127.0.0.1"),
		ListenPort: getEnvInt("LISTEN_PORT", 9240),

		UsageCollectorInterval:   getEnvInt("USAGE_COLLECTOR_INTERVAL", 30),
		BucketsCollectorInterval: getEnvInt("BUCKETS_COLLECTOR_INTERVAL", 300),
		UsersCollectorInterval:   getEnvInt("USERS_COLLECTOR_INTERVAL", 3600),

		RGWConnectionTimeout: getEnvInt("RGW_CONNECTION_TIMEOUT", 10),
		StartDelay:           getEnvInt("START_DELAY", 30),
		Insecure:             getEnvBool("INSECURE", false),
		SkipWithoutBucket:    getEnvBool("SKIP_WITHOUT_BUCKET", false),

		UsersCollectorEnable: getEnvBool("USERS_COLLECTOR_ENABLE", false),
	}

	var configFile string
	flag.StringVar(&configFile, "c", "", "config file")
	flag.Parse()

	if configFile == "" {
		// Work with ENV only
		return config, nil
	}

	s, err := os.Stat(configFile)
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, fmt.Errorf("'%s' is a directory, not a normal file", configFile)
	}

	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
