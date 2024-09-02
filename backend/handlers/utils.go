package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBoolQueryParam(c *gin.Context, key string, defValue bool) (bool, error) {
	param := c.Query(key)
	if param == "" {
		return false, nil
	}
	return strconv.ParseBool(param)
}
