package macro

import (
	"os"
	"path/filepath"
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
