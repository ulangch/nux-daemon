package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ulangch/nas_desktop_app/backend/models"
)

func CreatePrivateSpaceHandler(c *gin.Context) {
	var request struct {
		Password string `json:"password"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	space, err := models.CreatePrivateSpace(request.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "sid": space.Sid, "file": space.File, "token": space.Token})
	}
}

func GetPrivateSpaceHandler(c *gin.Context) {
	var request struct {
		Sid      string `json:"sid"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	space, err := models.GetPrivateSpace(request.Sid, request.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "sid": space.Sid, "file": space.File, "token": space.Token})
	}
}
