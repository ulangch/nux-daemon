package handlers

import (
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
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
	disks, err := models.GetDeviceDisks()
	if err != nil || len(disks) <= 0 {
		c.JSON(http.StatusOK, gin.H{"status_code": C_DISK_NOT_EXIST, "status_message": "no device disk"})
		return
	}
	var cpuUsage int
	if cpuPercents, err := cpu.Percent(0, false); err == nil {
		log.Printf("GetServiceInfo, cpuUsage=%f", cpuPercents[0])
		cpuUsage = int(math.Floor(cpuPercents[0]))
	}
	var memUsage int
	if memInfo, err := mem.VirtualMemory(); err == nil {
		log.Printf("GetServiceInfo, memUsage=%f", memInfo.UsedPercent)
		memUsage = int(math.Floor(memInfo.UsedPercent))
	}
	var url string
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				url = models.GenServiceUrl(ipNet.IP.String())
			}
		}
	}
	var qrcode string
	if url != "" {
		qcPath := filepath.Join(disks[0].Path, ".system")
		os.MkdirAll(qcPath, os.ModePerm)
		qcPath = filepath.Join(qcPath, "qc_device.png")
		qcData := models.GenQrCodeData(nid, url)
		if models.GenQrCodeFile(qcData, qcPath) == nil {
			qrcode = qcPath
		}
	}
	c.JSON(http.StatusOK, gin.H{"status_code": C_SUCCESS, "status_message": M_SUCCESS, "info": ServiceInfo{
		Name:   name,
		ID:     nid,
		CPU:    cpuUsage,
		Memory: memUsage,
		Url:    url,
		QRCode: qrcode,
		Disk:   disks[0].Path,
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
