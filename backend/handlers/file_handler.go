package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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
	file, err := models.GetFileInfo(decodePath, false)
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
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
		return
	}
	path := c.PostForm("path")
	if path == "" {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": "path invalid"})
		return
	}
	md5 := c.PostForm("md5")
	if md5 == "" {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": "md5 invalid"})
		return
	}
	ckNumber, err := strconv.ParseInt(c.PostForm("ck_number"), 10, 64)
	if err != nil || ckNumber <= 0 {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": "ck number invalid"})
	}
	ckSize, err := strconv.ParseInt(c.PostForm("ck_size"), 10, 64)
	if err != nil || ckSize <= 0 {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": "ck size invalid"})
	}
	fileDirPath := filepath.Dir(path)
	ckDirPath := filepath.Join(fileDirPath, ".uploads", md5)
	if err := os.MkdirAll(ckDirPath, 0777); err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": "failed to create upload dir"})
		return
	}
	ckFilePath := filepath.Join(ckDirPath, fmt.Sprintf("%d%s", ckNumber, ".ck"))
	// Save ck
	if err := c.SaveUploadedFile(file, ckFilePath); err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": "failed to save"})
		return
	}
	if ckNumber == ckSize {
		// Merge ck
		uploadFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": "failed to merge, path invalid"})
			return
		}
		defer uploadFile.Close()
		for i := 1; i <= int(ckSize); i++ {
			ckFile, err := os.Open(filepath.Join(ckDirPath, fmt.Sprintf("%d%s", i, ".ck")))
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": fmt.Sprintf("failed to merge, %d.ck invalid", i)})
				return
			}
			_, err = io.Copy(uploadFile, ckFile)
			ckFile.Close()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": fmt.Sprintf("failed to merge %d.ck", i)})
				return
			}
		}
		os.RemoveAll(ckDirPath)
	}
	c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
}

func GetUploadInfoHandler(c *gin.Context) {
	path, err := url.QueryUnescape(c.Query("path"))
	if err != nil || path == "" {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": "path invalid"})
		return
	}
	md5, err := url.QueryUnescape(c.Query("md5"))
	if err != nil || md5 == "" {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": "md5 invalid"})
		return
	}
	var file any = nil
	if fi, err := models.GetFileInfo(path, true); err == nil {
		file = fi
	}
	ckDirPath := filepath.Join(filepath.Dir(path), ".uploads", md5)
	var ckNumber = 0
	if entries, err := os.ReadDir(ckDirPath); err == nil {
		cks := make([]int, len(entries))
		ckNumber = len(entries)
		for _, entry := range entries {
			if ck, err := strconv.ParseInt(strings.TrimSuffix(entry.Name(), ".ck"), 10, 64); err == nil {
				if ck > 0 && ck <= int64(len(entries)) {
					cks[ck-1] = int(ck)
				}
			}
		}
		for index, ck := range cks {
			if ck <= 0 {
				ckNumber = index
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "file": file, "ck_number": ckNumber})
}

func UploadFileHandler2(c *gin.Context) {
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

	// Read the file ck from the request body
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
		return
	}

	// Upload the file ck
	err = models.WriteFileChunk(decodePath, offset, data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
}
