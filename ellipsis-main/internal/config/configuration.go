package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configuration struct {
	ServerPort    string `mapstructure:"SERVER_PORT"`
	SecretKey     string `mapstructure:"SECRET_KEY"`
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	PGHost        string `mapstructure:"PG_HOST"`
	PGPort        string `mapstructure:"PG_PORT"`
	PGDatabase    string `mapstructure:"PG_DATABASE"`
	PGUser        string `mapstructure:"PG_USER"`
	PGPassword    string `mapstructure:"PG_PASSWORD"`
	PGSSLMode     string `mapstructure:"PG_SSL_MODE"`
	AdminID       string `mapstructure:"ADMIN_ID"`
	AdminPassword string `mapstructure:"ADMIN_PASSWORD"`
}

// Setup initialize configuration
var (
	Config *Configuration
)

func Setup() {
	var configuration *Configuration
	logger := log.New()

	viper.SetConfigName("sample")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Overwrite file env's from environment
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		logger.Fatalf("Unable to decode into struct, %v", err)
	}

	if port := os.Getenv("PORT"); port != "" {
		configuration.ServerPort = port
	}

	Config = configuration
	logger.Info("configurations loading successfully")
}

// GetConfig helps you to get configuration data
func GetConfig() *Configuration {
	return Config
}
