package config

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

// ServerConfig holds basic server configuration.
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// InitialAdminConfig holds initial admin username and password.
type InitialAdminConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// AuthConfig holds authentication related configuration.
type AuthConfig struct {
	JWTSecret     string             `mapstructure:"jwt_secret"`
	InitialAdmin  InitialAdminConfig `mapstructure:"initial_admin"`
	TokenLifetime int64              `mapstructure:"token_lifetime_seconds"`
}

// AppConfig is the root configuration structure.
type AppConfig struct {
	Server ServerConfig `mapstructure:"server"`
	Auth   AuthConfig   `mapstructure:"auth"`
}

var (
	appConfig AppConfig
	once      sync.Once
)

// LoadConfig loads application configuration from YAML using Viper.
func LoadConfig() AppConfig {
	once.Do(func() {
		configPath := os.Getenv("FNCHATBOT_CONFIG")
		if configPath == "" {
			configPath = "config.yaml"
		}

		v := viper.New()
		v.SetConfigFile(configPath)
		v.SetConfigType("yaml")

		// Reasonable defaults for development.
		v.SetDefault("server.port", "8080")
		v.SetDefault("auth.jwt_secret", "change-me-in-config")
		v.SetDefault("auth.token_lifetime_seconds", 86400)

		if err := v.ReadInConfig(); err != nil {
			log.Printf("Config: unable to read config file %s, using defaults: %v", configPath, err)
		}

		if err := v.Unmarshal(&appConfig); err != nil {
			log.Printf("Config: failed to unmarshal config: %v", err)
		}
	})

	return appConfig
}

// GetConfig returns the loaded application configuration.
func GetConfig() AppConfig {
	// Ensure config is at least initialized with defaults.
	if (appConfig == AppConfig{}) {
		return LoadConfig()
	}
	return appConfig
}
