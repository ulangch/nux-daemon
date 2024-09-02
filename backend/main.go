package main

import (
	// "github.com/gin-gonic/gin"
	"log"

	"github.com/ulangch/nas_desktop_app/backend/config"
	"github.com/ulangch/nas_desktop_app/backend/models"
	"github.com/ulangch/nas_desktop_app/backend/routers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Set log level
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting server on port %s with log level %s", config.AppConfig.ServerPort, config.AppConfig.LogLevel)

	// Setup database
	db, err := gorm.Open(sqlite.Open(config.AppConfig.DBPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.KeyValue{})
	models.InitializeKVStore(db)

	// Setup device
	models.InitializeDeviceID()

	// Start the server
	r := routers.SetupRouter()
	r.Run(":" + config.AppConfig.ServerPort)
}
