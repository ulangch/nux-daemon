package client_handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/ulangch/nas_desktop_app/backend/core"
	"github.com/ulangch/nas_desktop_app/backend/utils"
)

func ListFilesHandler(c *gin.Context) {
	path, err := url.QueryUnescape(c.Query("path"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": utils.HTTP_CODE_FAILED, "status_message": err.Error()})
	}
	files, err := core.UnixRelativeListFiles(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": utils.HTTP_CODE_FAILED, "status_message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status_code": utils.HTTP_CODE_SUCCESS, "status_message": utils.HTTP_MSG_SUCCESS, "files": files})
}
