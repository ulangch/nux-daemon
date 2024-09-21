//go:build windows
// +build windows

package macro

import (
	"log"
	"path/filepath"
	"syscall"
	"unsafe"
)

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

func EncodeFilePath(unixPath string) string {
	return filepath.ToSlash(unixPath)
}

func DecodeFilePath(unixPath string) string {
	return filepath.FromSlash(unixPath)
}
