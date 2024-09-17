package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ulangch/nas_desktop_app/backend/models"
)

func AddRecentOpenHandler(c *gin.Context) {
	path, err := GetQueryRealPath(c, "path")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	models.AddRecentOpenFiles([]string{path})
	c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
}

func ListRecentOpenHandler(c *gin.Context) {
	files, err := models.ListRecentOpenFiles()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "files": files})
	}
}

func ListRecentAddHandler(c *gin.Context) {
	files, err := models.ListRecentAddFiles()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "files": files})
	}
}

func ListRecentDeleteHandler(c *gin.Context) {
	files, err := models.ListRecentDeleteFiles()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "files": files})
	}
}

func RecoverRecentDeleteHandler(c *gin.Context) {
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
	if err := models.RecoverRecentDeleteFiles(paths); err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
	}
}
