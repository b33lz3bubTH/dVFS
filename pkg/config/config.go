package config

import (
	"os"
	"strconv"

	"github.com/google/uuid"
)

// Config holds all configuration for the storage node
type Config struct {
	Port        int
	StoragePath string
	InstanceID  string
}

// Load reads configuration from environment variables with defaults
func Load() *Config {
	cfg := &Config{
		Port:        8080,
		StoragePath: "/tmp/storage_data",
	}

	// Override with environment variables if set
	if portStr := os.Getenv("PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.Port = port
		}
	}

	if storagePath := os.Getenv("STORAGE_PATH"); storagePath != "" {
		cfg.StoragePath = storagePath
	}

	// Load instance ID from environment or generate a new one
	if instanceID := os.Getenv("INSTANCE_ID"); instanceID != "" {
		cfg.InstanceID = instanceID
	} else {
		cfg.InstanceID = uuid.New().String()
	}

	return cfg
}
