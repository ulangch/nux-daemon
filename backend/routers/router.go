package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ulangch/nas_desktop_app/backend/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Define routes
	r.GET("/device/info", handlers.GetDeviceInfoHandler)
	r.POST("/device/update_name", handlers.UpdateDeviceNameHandler)
	r.POST("/device/add_disk", handlers.AddDiskPath)
	r.POST("/device/remove_disk", handlers.RemoveDiskPath)
	r.GET("/file/list", handlers.ListFilesHandler)
	r.GET("/file/info", handlers.GetFileInfoHandler)
	r.GET("/file/stream", handlers.ReadFileHandler)
	r.POST("/file/create", handlers.CreateFileHandler)
	r.POST("/file/mkdir", handlers.CreateDirectoryHandler)
	r.POST("/file/delete", handlers.DeleteFileHandler)
	r.POST("/file/move", handlers.MoveFileHandler)
	r.POST("/file/upload", handlers.UploadFileHandler)
	r.GET("/file/upload_info", handlers.GetUploadInfoHandler)
	r.GET("/gallery/list_files", handlers.ListGalleryFilesHandler)
	r.GET("/gallery/get_dir", handlers.GetGalleryDirHandler)
	r.POST("/gallery/update_dir", handlers.UpdateGalleryDirHandler)
	return r
}
