package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ulangch/nas_desktop_app/backend/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Define routes
	r.POST("/file/create", handlers.CreateFileHandler)
    r.POST("/file/mkdir", handlers.CreateDirectoryHandler)
    r.GET("/file/info", handlers.GetFileInfoHandler)
	r.GET("/file/read", handlers.ReadFileHandler)
	r.DELETE("/file/delete", handlers.DeleteFileHandler)
	r.GET("/file/list", handlers.ListFilesHandler)

	return r
}
