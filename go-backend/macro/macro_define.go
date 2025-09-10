package macro

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type DiskUsage struct {
	Total int64
	Free  int64
	Used  int64
}

func GetSystemDirPath() string {
	path := ".system"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)

	}
	if absolutePath, err := filepath.Abs(path); err == nil {
		return absolutePath
	} else {
		return path
	}
}

func GetSystemMemory() (int64, int64) {
	v, _ := mem.VirtualMemory()
	return int64(v.Total), int64(v.Available)
}

func GetSystemCpuRate() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return math.Round(percent[0]*10) / 10
}

func GetSystemCpuTemperature() float64 {
	temps, err := host.SensorsTemperatures()
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	for _, t := range temps {
		if t.SensorKey == "cpu-thermal" || t.SensorKey == "coretemp_packageid0" {
			return math.Round(t.Temperature*10) / 10
		}
	}
	return 0
}
