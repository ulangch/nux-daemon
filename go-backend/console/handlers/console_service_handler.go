package console_handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/ulangch/nas_desktop_app/backend/core"
	"github.com/ulangch/nas_desktop_app/backend/utils"
)

func GetServiceInfoHandler(c *gin.Context) {
	nuxDevice, err := core.GetDevice()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": utils.HTTP_CODE_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": utils.HTTP_CODE_SUCCESS, "status_message": utils.HTTP_MSG_SUCCESS, "nux_device": nuxDevice})
	}
}

func AddDiskHandler(c *gin.Context) {
	path, err := url.QueryUnescape(c.Query("path"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": utils.HTTP_CODE_FAILED, "status_message": err.Error()})
		return
	}
	err = core.AddDiskByPlatformPath(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": utils.HTTP_CODE_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": utils.HTTP_CODE_SUCCESS, "status_message": utils.HTTP_MSG_SUCCESS})
	}
}

func RemoveDiskHandler(c *gin.Context) {
	path, err := url.QueryUnescape(c.Query("path"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": utils.HTTP_CODE_FAILED, "status_message": err.Error()})
		return
	}
	err = core.RemoveDiskByPlatformPath(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": utils.HTTP_CODE_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": utils.HTTP_CODE_SUCCESS, "status_message": utils.HTTP_MSG_SUCCESS})
	}
}
