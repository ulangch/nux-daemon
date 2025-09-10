package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ulangch/nas_desktop_app/backend/utils"
)

/**
 * unixRelativePath: /{DiskId}/backup
 */
func UnixRelativeListFiles(unixRelativePath string) ([]NuxFile, error) {
	unixAbsolutePath, err := UnixRelative2Absolute(unixRelativePath)
	if err == nil {
		return UnixAbsoluteListFiles(unixAbsolutePath)
	} else {
		return nil, err
	}
}

/**
 * unixAbsolutePath:
 * unix: /Users/uki/nux/backup
 * win: /c/uki/nux/backup
 */
func UnixAbsoluteListFiles(unixAbsolutePath string) ([]NuxFile, error) {
	if platformAbsolutePath := utils.ToPlatformPath(unixAbsolutePath); platformAbsolutePath != "" {
		return PlatformListFiles(platformAbsolutePath)
	} else {
		err := fmt.Errorf("core_file UnixAbsoluteListFiles, ToPlatformPath failed, unixAbsolutePath=%s", unixAbsolutePath)
		return nil, err
	}
}

/**
 * platformPath:
 * unix: /Users/uki/nux/backup
 * win: C:\uki\nux\backup
 */
func PlatformListFiles(platformPath string) ([]NuxFile, error) {
	disk, err := GetDiskByAbsolutePath(platformPath)
	if err != nil || disk.ID == "" {
		return nil, fmt.Errorf("core_file PlatformListFiles, GetDiskByAbsolutePath failed, platformPath=%s", platformPath)
	}
	entries, err := os.ReadDir(platformPath)
	if err != nil {
		return nil, err
	}
	var nuxFiles []NuxFile
	for _, entry := range entries {
		if !utils.IsValidFileName(entry.Name()) {
			continue
		}
		fileInfo, err := entry.Info()
		if err != nil {
			continue
		}
		filepath := filepath.Join(platformPath, entry.Name())
		nuxFile, _ := PackNuxFile(disk, filepath, fileInfo)
		nuxFiles = append(nuxFiles, nuxFile)
	}
	return nuxFiles, nil
}

func PackNuxFile(disk NuxDisk, platformPath string, fileInfo os.FileInfo) (NuxFile, error) {
	unixAbsolutePath := utils.ToUnixPath(platformPath)
	unixRelativePath, _ := UnixAbsolute2Relative2(disk, unixAbsolutePath)
	nuxFile := NuxFile{
		Name:         fileInfo.Name(),
		UnixPath:     unixRelativePath,
		UnixAbsolute: unixAbsolutePath,
		Size:         fileInfo.Size(),
		UpdateTime:   fileInfo.ModTime().UnixMilli(),
		IsDir:        fileInfo.IsDir(),
		MD5:          utils.GetFileMD5(platformPath),
		Thumbnail:    GetImageThumbnail(platformPath),
		IsCollected:  IsFileCollected(unixRelativePath),
	}
	return nuxFile, nil
}

/**
 * unix: /{DiskId}/backup => /Users/uki/nux/backup
 * win:  /{DiskId}/backup => /c/uki/nux/backup
 */
func UnixRelative2Absolute(unixRelativePath string) (string, error) {
	diskId := utils.UnixRootDir(unixRelativePath)
	if diskId == "" {
		return "", fmt.Errorf("core_file UnixRelative2Absolute, UnixRootDir failed, unixRelativePath=%s", unixRelativePath)
	}
	disk, err := GetDiskById(diskId)
	if err != nil || disk.ID == "" {
		return "", fmt.Errorf("core_file UnixRelative2Absolute, no matched disk, unixRelativePath=%s", unixRelativePath)
	}
	diskUnixAbsolute := disk.UnixAbsolute // /Users/uki/nux(unix), /c/uki/nux(win)
	diskUnixAbsolute = strings.TrimPrefix(diskUnixAbsolute, "/")
	diskUnixAbsolute = strings.TrimSuffix(diskUnixAbsolute, "/")
	absolutePath := strings.ReplaceAll(unixRelativePath, diskId, diskUnixAbsolute)
	return absolutePath, nil
}

/**
 * unix: /User/uki/nux/backup => /{DiskId}/backup
 * win:  /c/uki/nux/backup    => /{DiskId}/backup
 */
func UnixAbsolute2Relative(unixAbsolutePath string) (string, error) {
	disk, err := GetDiskByAbsolutePath(unixAbsolutePath)
	if err == nil {
		return UnixAbsolute2Relative2(disk, unixAbsolutePath)
	} else {
		return "", err
	}
}

/**
 * unix: /User/uki/nux/backup => /{DiskId}/backup
 * win:  /c/uki/nux/backup    => /{DiskId}/backup
 */
func UnixAbsolute2Relative2(disk NuxDisk, unixAbsolutePath string) (string, error) {
	diskUnixAbsolute := disk.UnixAbsolute
	diskUnixAbsolute = strings.TrimPrefix(diskUnixAbsolute, "/")
	diskUnixAbsolute = strings.TrimSuffix(diskUnixAbsolute, "/")
	relativePath := strings.ReplaceAll(unixAbsolutePath, diskUnixAbsolute, disk.ID)
	return relativePath, nil
}

/**
 * unix: /User/uki/nux/backup => /User/uki/nux/backup
 * win:  /c/uki/nux/backup    => C:\uki\nux\backup
 */
func UnixAbsolute2Platform(unixAbsolutePath string) (string, error) {
	return utils.ToPlatformPath(unixAbsolutePath), nil
}
