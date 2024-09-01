package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ulangch/nas_desktop_app/backend/models"
)

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
		return
	}

	c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "file": file})
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
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "file": file})
}

func GetFileInfoHandler(c *gin.Context) {
	path := c.Query("path")
	decodePath, err := url.QueryUnescape(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	file, err := models.GetFileInfo(decodePath)
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
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
}

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
