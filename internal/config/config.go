package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	defaultServerPort         = 8080
	defaultJWTExpirationHours = 72
)

// Config represents an application configuration.
type Config struct {
	// The server port. Defaults to 8080
	ServerPort int `mapstructure:"server_port"`
	// The data source name (DSN) for connecting to the database. required.
	DSN string `mapstructure:"dsn"`
	// JWT signing key. required.
	JWTSigningKey string `mapstructure:"jwt_signing_key"`
	// JWT expiration in hours. Defaults to 72 hours (3 days)
	JWTExpiration int `mapstructure:"jwt_expiration"`
}

// Load returns an application configuration which is populated from the given configuration file and environment variables.
func Load(filename string) (*Config, error) {
	v := viper.New()
	// Default config
	v.SetDefault("server_port", defaultServerPort)
	v.SetDefault("jwt_expiration", defaultJWTExpirationHours)

	// load from YAML config file
	v.SetConfigName(filename)
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.SetEnvPrefix("app")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %s", err)
		}
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}

	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %s", err)
	}
	return &c, nil
}
