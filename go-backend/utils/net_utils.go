package utils

import (
	"fmt"
	"net"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetLocalIPv4() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		// 跳过未启用或回环接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 过滤IPv4并且不是回环地址
			if ip != nil && ip.To4() != nil && !ip.IsLoopback() {
				return ip.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no active IPv4 address found")
}

func GetBoolQueryParam(c *gin.Context, key string, defValue bool) (bool, error) {
	param := c.Query(key)
	if param == "" {
		return false, nil
	}
	return strconv.ParseBool(param)
}
