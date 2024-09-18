package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ulangch/nas_desktop_app/backend/macro"
)

// Config holds the application configuration
type Config struct {
	ServerPort string `json:"server_port"`
	LogLevel   string `json:"log_level"`
	DBPath     string `json:"db_path"`
}

// AppConfig is the global configuration instance
var AppConfig Config

// LoadConfig loads the configuration from a file or environment variables
func LoadConfig() {
	// Try to load configuration from file
	file, err := os.Open("config.json")
	if err != nil {
		log.Printf("Error opening config file: %v. Falling back to environment variables.", err)
		loadConfigFromEnv()
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		log.Printf("Error decoding config file: %v. Falling back to environment variables.", err)
		loadConfigFromEnv()
	}
}

// loadConfigFromEnv loads configuration from environment variables
func loadConfigFromEnv() {
	AppConfig.ServerPort = getEnv("SERVER_PORT", "8080")
	AppConfig.LogLevel = getEnv("LOG_LEVEL", "info")
	// AppConfig.DBPath = getEnv("DB_PATH", "my-nas-app.db")
	AppConfig.DBPath = macro.GetDatabasePath()
}

// getEnv gets the value of an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
