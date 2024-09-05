package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ulangch/nas_desktop_app/backend/models"
)

// ListFilesHandler handles listing all files in a directory
func ListFilesHandler(c *gin.Context) {
	dir := c.Query("path")
	decodedDir, err := url.QueryUnescape(dir)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	files, err := models.ListFiles(decodedDir)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_message": "", "files": files})
}

func GetFileInfoHandler(c *gin.Context) {
	path := c.Query("path")
	decodePath, err := url.QueryUnescape(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	file, err := models.GetFileInfo(decodePath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "file": nil})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "file": file})
}

// ReadFileHandler handles reading a file
func ReadFileHandler(c *gin.Context) {
	path := c.Query("path")
	decodePath, err := url.QueryUnescape(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	// 打开文件
	file, err := os.Open(decodePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
		return
	}
	// 确保在函数结束时关闭文件
	defer file.Close()
	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
		return
	}
	// 设置下载文件的头信息
	c.Header("Content-Disposition", "attachment; filename="+fileInfo.Name())
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprint(fileInfo.Size()))
	// 流式传输文件内容到响应
	http.ServeFile(c.Writer, c.Request, path)
}

// CreateFileHandler handles the creation of a new file
func CreateFileHandler(c *gin.Context) {
	path := c.Query("path")
	decodePath, err := url.QueryUnescape(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	file, err := models.CreateFile(decodePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "file": file})
	}
}

func CreateDirectoryHandler(c *gin.Context) {
	path := c.Query("path")
	decodePath, err := url.QueryUnescape(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	file, err := models.CreateDirectory(decodePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "file": file})
	}
}

// DeleteFileHandler handles deleting a file
func DeleteFileHandler(c *gin.Context) {
	path := c.Query("path")
	decodePath, err := url.QueryUnescape(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	err = models.DeleteFile(decodePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
	}
}

func MoveFileHandler(c *gin.Context) {
	oldPath, err := url.QueryUnescape(c.Query("old_path"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	newPath, err := url.QueryUnescape(c.Query("new_path"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	err = models.MoveFile(oldPath, newPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
	}
}

func UploadFileHandler(c *gin.Context) {
	path := c.Query("path")
	decodePath, err := url.QueryUnescape(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	offsetStr := c.Query("offset")
	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": "Invalid offset"})
		return
	}

	// Read the file chunk from the request body
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
		return
	}

	// Upload the file chunk
	err = models.WriteFileChunk(decodePath, offset, data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
}

func UploadFileChunkHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
		return
	}

	// 创建保存文件的目录
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// 保存文件
	filePath := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file_path": filePath})
}

func MergeFileChunkHandler(c *gin.Context) {

}
