package main

import (
	// "github.com/gin-gonic/gin"
	"log"

	"github.com/ulangch/nas_desktop_app/backend/config"
	"github.com/ulangch/nas_desktop_app/backend/routers"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Set log level
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting server on port %s with log level %s", config.AppConfig.ServerPort, config.AppConfig.LogLevel)

	// Start the server
	r := routers.SetupRouter()
	r.Run(":" + config.AppConfig.ServerPort)
}
