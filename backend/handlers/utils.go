package handlers

import (
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ulangch/nas_desktop_app/backend/models"
)

func GetBoolQueryParam(c *gin.Context, key string, defValue bool) (bool, error) {
	param := c.Query(key)
	if param == "" {
		return false, nil
	}
	return strconv.ParseBool(param)
}

func GetQueryRealPath(c *gin.Context, key string) (string, error) {
	clientPath, err := url.QueryUnescape(c.Query(key))
	if err != nil {
		return "", err
	}
	return models.GetRealPath(clientPath)
}
