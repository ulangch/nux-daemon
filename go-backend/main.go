package main

import (
	// "github.com/gin-gonic/gin"
	"log"
	"os"

	"github.com/ulangch/nas_desktop_app/backend/config"
	"github.com/ulangch/nas_desktop_app/backend/models"
	"github.com/ulangch/nas_desktop_app/backend/routers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Set log level
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("starting server on port %s with log level %s", config.AppConfig.ServerPort, config.AppConfig.LogLevel)

	dbPath := config.AppConfig.DBPath
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatalf("failed to create database file: %v", err)
		}
		file.Close()
	}

	// Setup database
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.KeyValue{}, &models.Collection{}, &models.RecentOpenFile{}, &models.RecentAddFile{}, &models.RecentDeleteFile{})
	models.InitializeKVStore(db)
	models.InitializeColStore(db)
	models.InitializeRecentDB(db)

	// Setup device
	models.InitializeDeviceID()

	// Start the server
	r := routers.SetupRouter()
	r.Run(":" + config.AppConfig.ServerPort)
}
