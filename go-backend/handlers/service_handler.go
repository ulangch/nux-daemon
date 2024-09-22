package handlers

import (
	"math"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/ulangch/nas_desktop_app/backend/macro"
	"github.com/ulangch/nas_desktop_app/backend/models"
)

type ServiceInfo struct {
	Name   string `json:"name"`
	ID     string `json:"id"`
	CPU    int    `json:"cpu"`
	Memory int    `json:"memory"`
	Url    string `json:"url"`
	QRCode string `json:"qrcode"`
	Disk   string `json:"disk"`
}

func GetServiceInfo(c *gin.Context) {
	nid, err := models.GetKV(models.KV_KEY_DEVICE_ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": "no device id"})
		return
	}
	name, err := models.GetKV(models.KV_KEY_DEVICE_NAME)
	if err != nil {
		name = models.DEFAULT_DEVICE_NAME
	}
	var cpuUsage int
	if cpuPercents, err := cpu.Percent(0, false); err == nil {
		cpuUsage = int(math.Min(math.Ceil(cpuPercents[0]), 100))
	}
	var memUsage int
	if memInfo, err := mem.VirtualMemory(); err == nil {
		memUsage = int(math.Floor(memInfo.UsedPercent))
	}
	var url string
	if ipv4 := models.GetLocalIPV4(); ipv4 != "" {
		url = models.GenServiceUrl(ipv4)
	}
	var qrcode string
	if url != "" {
		qcPath := filepath.Join(macro.GetSystemDirPath(), "qc_device.png")
		qcData := models.GenQrCodeData(nid, url)
		if models.GenQrCodeFile(qcData, qcPath) == nil {
			qrcode = qcPath
		}
	}
	var diskPath string
	if disks, err := models.GetDeviceDisks(); err == nil && len(disks) > 0 {
		diskPath = disks[0].Path
	}
	c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "info": ServiceInfo{
		Name:   name,
		ID:     nid,
		CPU:    cpuUsage,
		Memory: memUsage,
		Url:    url,
		QRCode: qrcode,
		Disk:   diskPath,
	}})
}

func UpdateDeviceName(c *gin.Context) {
	name, err := url.QueryUnescape(c.Query("name"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	err = models.UpdateDeviceName(name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
	}
}

func UpdateDiskPath(c *gin.Context) {
	path, err := url.QueryUnescape(c.Query("path"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_INVALID_PARAM, "status_message": err.Error()})
		return
	}
	err = models.UpdateDiskPath(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status_code": C_REQUEST_FAILED, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS})
	}
}
