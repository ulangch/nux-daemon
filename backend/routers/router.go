package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ulangch/nas_desktop_app/backend/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Define routes
	r.GET("/device/info", handlers.GetDeviceInfoHandler)
	r.POST("/device/set_name", handlers.UpdateDeviceNameHandler)
	r.POST("/device/add_disk", handlers.AddDiskPath)
	r.POST("/device/remove_disk", handlers.RemoveDiskPath)

	r.GET("/file/list", handlers.ListFilesHandler)
	r.GET("/file/info", handlers.GetFileInfoHandler)
	r.GET("/file/stream", handlers.StreamFileHandler)
	r.GET("/file/stream_seek", handlers.StreamSeekFileHandler)
	r.GET("/file/stream_thumbnail", handlers.StreamThumbnailFileHandler)
	r.POST("/file/create", handlers.CreateFileHandler)
	r.POST("/file/mkdir", handlers.CreateDirectoryHandler)
	r.POST("/file/delete", handlers.DeleteFileHandler)
	r.POST("/file/delete_batch", handlers.BatchDeleteFileHandler)
	r.POST("/file/move", handlers.MoveFileHandler)
	r.POST("/file/move_batch", handlers.BatchMoveFilesHandler)
	r.POST("/file/upload", handlers.UploadFileHandler)
	r.GET("/file/upload_info", handlers.GetUploadInfoHandler)

	r.GET("/client/gallery/get_dir", handlers.GetGalleryDirHandler)
	r.POST("/client/gallery/set_dir", handlers.UpdateGalleryDirHandler)
	r.GET("/client/gallery/list_files", handlers.ListGalleryFilesHandler)

	r.GET("/client/upload/get_dir", handlers.GetUploadDirHandler)
	r.POST("/client/upload/set_dir", handlers.SetUploadDirHandler)

	r.GET("/client/bucket/get_dirs", handlers.GetBucketDirsHandler)
	r.POST("/client/bucket/delete_dirs", handlers.DeleteBucketDirsHandler)
	r.POST("/client/bucket/add_dirs", handlers.AddBucketDirsHandler)
	r.GET("/client/bucket/list_files", handlers.ListBucketFilesHandler)

	r.POST("/client/collect/add", handlers.CollectFilesHandler)
	r.POST("/client/collect/delete", handlers.UnCollectFilesHandler)
	r.GET("/client/collect/list_files", handlers.ListCollectFilesHandler)

	r.POST("/client/recent/add_open_files", handlers.AddRecentOpenHandler)
	r.GET("/client/recent/list_open_files", handlers.ListRecentOpenHandler)
	r.GET("/client/recent/list_upload_files", handlers.ListRecentAddHandler)
	r.GET("/client/recent/list_delete_files", handlers.ListRecentDeleteHandler)
	r.POST("/client/recent/recover_delete_files", handlers.RecoverRecentDeleteHandler)
	return r
}
