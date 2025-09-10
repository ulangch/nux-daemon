//go:build windows
// +build windows

package macro

import (
	"log"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

func IsWin() bool {
	return true
}

func GetDatabasePath() string {
	return filepath.Join(GetSystemDirPath(), "nas-daemon-windows.db")
}

func GetDiskUsage(path string) (DiskUsage, error) {
	kernel32, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return DiskUsage{}, err
	}
	defer syscall.FreeLibrary(kernel32)
	GetDiskFreeSpaceEx, err := syscall.GetProcAddress(syscall.Handle(kernel32), "GetDiskFreeSpaceExW")

	if err != nil {
		return DiskUsage{}, err
	}

	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)

	diskDirUTF16Ptr, err := syscall.UTF16PtrFromString("F:")
	if err != nil {
		return DiskUsage{}, err
	}

	syscall.Syscall6(uintptr(GetDiskFreeSpaceEx), 4,
		uintptr(unsafe.Pointer(diskDirUTF16Ptr)),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)), 0, 0)

	log.Printf("Available  %dmb", lpFreeBytesAvailable/1024/1024.0)
	log.Printf("Total      %dmb", lpTotalNumberOfBytes/1024/1024.0)
	log.Printf("Free       %dmb", lpTotalNumberOfFreeBytes/1024/1024.0)

	return DiskUsage{Total: lpTotalNumberOfBytes, Free: lpTotalNumberOfFreeBytes, Used: lpTotalNumberOfBytes - lpTotalNumberOfFreeBytes}, nil
}

func IsSameVolume(path1, path2 string) bool {
	vol1 := strings.ToUpper(filepath.VolumeName(path1))
	vol2 := strings.ToUpper(filepath.VolumeName(path2))
	return vol1 == vol2
}

func EncodeFilePath(unixPath string) string {
	return filepath.ToSlash(unixPath)
}

func DecodeFilePath(unixPath string) string {
	return filepath.FromSlash(unixPath)
}

func UnixToPlatformPath(unixPath string) string {
	p := unixPath
	// 如果是 /c/... 这种格式 → C:\...
	if strings.HasPrefix(p, "/") && len(p) > 2 && p[2] == '/' {
		drive := strings.ToUpper(string(p[1]))
		p = drive + ":" + p[2:]
	}
	// 替换 / 为 \
	p = strings.ReplaceAll(p, "/", `\`)
	return p
}
