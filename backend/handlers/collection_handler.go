package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ulangch/nas_desktop_app/backend/models"
)

func CollectFilesHandler(c *gin.Context) {
	var request struct {
		Paths []string `json:"paths"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	models.CollectFiles(request.Paths)
	c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
}

func UnCollectFilesHandler(c *gin.Context) {
	var request struct {
		Paths []string `json:"paths"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	models.UnCollectFiles(request.Paths)
	c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
}

func ListCollectFilesHandler(c *gin.Context) {
	files, err := models.ListCollectFiles()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "files": files})
	}
}
