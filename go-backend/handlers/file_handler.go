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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulangch/nas_desktop_app/backend/models"
)

//  /[disk]/[path]

// ListFilesHandler handles listing all files in a directory
func ListFilesHandler(c *gin.Context) {
	path, err := GetQueryRealPath(c, "path")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
	}
	files, err := models.ListFiles(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_message": M_SUCCESS, "files": files})
}

func GetFileInfoHandler(c *gin.Context) {
	path, err := GetQueryRealPath(c, "path")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
	}
	file, err := models.GetFileInfo(path)
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

func StreamFileHandler(c *gin.Context) {
	path, err := GetQueryRealPath(c, "path")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status_message": err.Error()})
		return
	}
	file, err := os.Open(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status_message": err.Error()})
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status_message": err.Error()})
		return
	}
	c.Header("Content-Disposition", "attachment; filename="+fileInfo.Name())
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprint(fileInfo.Size()))
	http.ServeFile(c.Writer, c.Request, path)
}

func StreamSeekFileHandler(c *gin.Context) {
	path, err := GetQueryRealPath(c, "path")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status_message": err.Error()})
		return
	}

	file, err := os.Open(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status_message": err.Error()})
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status_message": err.Error()})
		return
	}
	var start int64 = 0
	var end int64 = fileInfo.Size() - 1
	rangeHeader := c.GetHeader("Range")
	if rangeHeader != "" {
		ranges := strings.Split(rangeHeader, "=")
		if len(ranges) != 2 || ranges[0] != "bytes" {
			c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"status_message": "Invalid Range header"})
			return
		}
		rangeParts := strings.Split(ranges[1], "-")
		if len(rangeParts) != 2 {
			c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"status_message": "Invalid Range header"})
			return
		}
		start, err = strconv.ParseInt(rangeParts[0], 10, 64)
		if err != nil {
			c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"status_message": "Invalid Range header"})
			return
		}
		if rangeParts[1] != "" {
			end, err = strconv.ParseInt(rangeParts[1], 10, 64)
			if err != nil {
				c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"error": "Invalid Range header"})
				return
			}
		}
	}
	if start > end || start < 0 || end >= fileInfo.Size() {
		c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"error": "Invalid Range header"})
		return
	}
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileInfo.Size()))
	c.Header("Content-Length", fmt.Sprintf("%d", end-start+1))
	c.Status(http.StatusPartialContent)

	file.Seek(start, 0)
	io.CopyN(c.Writer, file, end-start+1)
}

func StreamThumbnailFileHandler(c *gin.Context) {
	path, err := GetQueryRealPath(c, "path")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status_message": err.Error()})
		return
	}
	md5, err := models.GetFileMD5(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status_message": err.Error()})
		return
	}
	disks, err := models.GetDeviceDisks()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status_message": err.Error()})
		return

	}
	thumbnailPath := filepath.Join(disks[0].Path, ".thumbnail", fmt.Sprintf("%s.JPEG", md5))
	thumbnail, err := os.Open(thumbnailPath)
	if err != nil {
		// TODO: Generate thumbnail
		c.JSON(http.StatusNotFound, gin.H{"status_message": err.Error()})
		return
	}
	defer thumbnail.Close()
	fileInfo, err := thumbnail.Stat()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status_message": err.Error()})
		return
	}
	c.Header("Content-Disposition", "attachment; filename="+fileInfo.Name())
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprint(fileInfo.Size()))
	http.ServeFile(c.Writer, c.Request, thumbnailPath)
}

func CreateFileHandler(c *gin.Context) {
	path, err := GetQueryRealPath(c, "path")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	file, err := models.CreateFile(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "file": file})
	}
}

func CreateDirectoryHandler(c *gin.Context) {
	path, err := GetQueryRealPath(c, "path")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	file, err := models.CreateDirectory(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "file": file})
	}
}

// DeleteFileHandler handles deleting a file
func DeleteFileHandler(c *gin.Context) {
	path, err := GetQueryRealPath(c, "path")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	err = models.DeleteFile(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
	}
}

func BatchDeleteFileHandler(c *gin.Context) {
	var request struct {
		Paths []string `json:"paths"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	paths, err := models.GetRealPaths(request.Paths)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	for _, path := range paths {
		if err := models.DeleteFile(path); err != nil {
			c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
}

func MoveFileHandler(c *gin.Context) {
	oldPath, err := GetQueryRealPath(c, "old_path")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	newPath, err := GetQueryRealPath(c, "new_path")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	file, err := models.MoveFile(oldPath, newPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "file": file})
	}
}

func BatchMoveFilesHandler(c *gin.Context) {
	var request struct {
		Paths []string `json:"paths"`
		Dir   string   `json:"dir"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	paths, err := models.GetRealPaths(request.Paths)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	dir, err := models.GetRealPath(request.Dir)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	if err := models.BatchMoveFiles(paths, dir); err != nil {
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
	path, err := models.GetRealPath(c.PostForm("path"))
	if err != nil {
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
		// Save file stat
		lastModified, err := strconv.ParseInt(c.PostForm("last_modified"), 10, 64)
		if err == nil && lastModified > 0 {
			os.Chtimes(path, time.UnixMilli(lastModified), time.UnixMilli(lastModified))
		}

		// Save thumbnail if have
		thumbnail, err := c.FormFile("thumb")
		disks, diskErr := models.GetDeviceDisks()
		if err == nil && diskErr == nil {
			thumbDirPath := filepath.Join(disks[0].Path, ".thumbnail")
			if err := os.MkdirAll(thumbDirPath, 0777); err == nil {
				thumbFilePath := filepath.Join(thumbDirPath, fmt.Sprintf("%s.JPEG", md5))
				c.SaveUploadedFile(thumbnail, thumbFilePath)
			}
		}
		// // Manage gallery
		// isGallery, err := strconv.ParseBool(c.PostForm("is_gallery"))

		// Add record
		models.AddRecentAddFiles([]string{path})
	}
	c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
}

func GetUploadInfoHandler(c *gin.Context) {
	path, err := GetQueryRealPath(c, "path")
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
	if fi, err := models.GetFileInfo(path); err == nil {
		file = fi
	}
	ckDirPath := filepath.Join(filepath.Dir(path), ".uploads", md5)
	var ckNumber = 0
	if entries, err := os.ReadDir(ckDirPath); err == nil {
		ckNumber = len(entries)
		for i := 1; i <= len(entries); i++ {
			ckName := fmt.Sprintf("%d.ck", i)
			hasCk := false
			for _, entry := range entries {
				if entry.Name() == ckName {
					hasCk = true
					break
				}
			}
			if !hasCk {
				ckNumber = i - 1
				break
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "file": file, "ck_number": ckNumber})
}
