package handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/ulangch/nas_desktop_app/backend/models"
)

func GetDeviceInfoHandler(c *gin.Context) {
	device, err := models.GetDeviceInfo()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "device": device})
	}
}

func UpdateDeviceNameHandler(c *gin.Context) {
	name := c.Query("name")
	decodeName, err := url.QueryUnescape(name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	err = models.UpdateDeviceName(decodeName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
	}
}

func AddDiskPath(c *gin.Context) {
	path := c.Query("path")
	decodePath, err := url.QueryUnescape(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	autoCreate, err := GetBoolQueryParam(c, "auto_create", false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	err = models.AddDiskPath(decodePath, autoCreate)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
	}
}

func RemoveDiskPath(c *gin.Context) {
	path := c.Query("path")
	decodePath, err := url.QueryUnescape(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	err = models.RemoveDiskPath(decodePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
	}
}
