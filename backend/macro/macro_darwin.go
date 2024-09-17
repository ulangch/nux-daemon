//go:build darwin
// +build darwin

package macro

import (
	"log"
	"syscall"
)

func GetDatabasePath() string {
	return "my-nas-app.db"
}

func GetDiskUsage(path string) (DiskUsage, error) {
	var stat syscall.Statfs_t
	var total uint64
	var free uint64
	err := syscall.Statfs(path, &stat)
	if err == nil {
		total = stat.Blocks * uint64(stat.Bsize)
		free = stat.Bfree * uint64(stat.Bsize)
		return DiskUsage{Total: int64(total), Free: int64(free), Used: int64(total) - int64(free)}, nil
	} else {
		log.Printf("GetDeviceInfo get volume failed: %s", err.Error())
		return DiskUsage{}, err
	}
}

func EncodeFilePath(unixPath string) string {
	return unixPath
}

func DecodeFilePath(unixPath string) string {
	return unixPath
}
