//go:build darwin
// +build darwin

package macro

import (
	"log"
	"os"
	"path/filepath"
	"syscall"
)

func IsWin() bool {
	return false
}

func GetDatabasePath() string {
	return filepath.Join(GetSystemDirPath(), "nas-daemon-darwin.db")
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

func IsSameVolume(path1, path2 string) bool {
	info1, err := os.Stat(path1)
	if err != nil {
		return false
	}
	info2, err := os.Stat(path2)
	if err != nil {
		return false
	}

	stat1 := info1.Sys().(*syscall.Stat_t)
	stat2 := info2.Sys().(*syscall.Stat_t)

	return stat1.Dev == stat2.Dev
}

func EncodeFilePath(unixPath string) string {
	return unixPath
}

func DecodeFilePath(unixPath string) string {
	return unixPath
}

func UnixToPlatformPath(unixPath string) string {
	return unixPath
}
