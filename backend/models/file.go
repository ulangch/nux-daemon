package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type File struct {
	Nid         string `json:"nid"`
	Name        string `json:"name"`
	Path        string
	Url         string `json:"url"`
	Size        int64  `json:"size"`
	UpdateTime  int64  `json:"update_time"`
	IsDir       bool   `json:"is_dir"`
	MD5         string `json:"md5"`
	Thumbnail   string `json:"thumbnail"`
	IsCollected bool   `json:"is_col"`
	GhostUrl    string `json:"ghost_url"` // only for recent delete
	FreeVolume  int64  `json:"free_volume"`
	TotalVolume int64  `json:"total_volume"`
}

// ListFiles lists all files in the specified directory
func ListFiles(dir string) ([]File, error) {
	var files []File
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	deviceID := GetDeviceID()
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if !FilterFile(entry.Name()) {
			continue
		}
		files = append(files, PackFileByInfo(filepath.Join(dir, info.Name()), info, deviceID))
	}
	return files, nil
}

func ListTypeFiles(dir string, filterType string, maxDepth int) ([]File, error) {
	var files []File
	if err := ListTypeFilesRecursive(&files, dir, filterType, 0, maxDepth); err != nil {
		return nil, err
	} else {
		return files, nil
	}
}

// List [filterType] files in [dir], depth start with 0, maxDepth should be equal or bigger than 0
func ListTypeFilesRecursive(result *[]File, dir string, filterType string, depth int, maxDepth int) error {
	if depth > maxDepth {
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}
	deviceID := GetDeviceID()
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() && entry.Name()[0] != '.' {
			if depth < maxDepth {
				ListTypeFilesRecursive(result, path, filterType, depth+1, maxDepth)
			}
			continue
		}
		if !FilterFile(entry.Name()) {
			continue
		}
		if !FilterFileByType(entry.Name(), filterType) {
			continue
		}
		*result = append(*result, PackFileByInfo(path, info, deviceID))
	}
	return nil
}

func PackFileByInfo(path string, info os.FileInfo, nid string) File {
	var md5 string
	if !info.IsDir() {
		md5, _ = GetFileMD5(path)
	}
	var thumbnail string
	if md5 != "" && (IsImage(info.Name()) || IsVideo(info.Name())) {
		thumbnail, _ = GetImageThumbnail(path, md5)
		if thumbnail != "" {
			thumbnail = fmt.Sprintf("nas://%s%s", nid, thumbnail)
		}
	}
	isCollected := IsCollected(path)
	return File{
		Nid:         nid,
		Name:        info.Name(),
		Path:        path,
		Url:         fmt.Sprintf("nas://%s%s", nid, path),
		Size:        info.Size(),
		UpdateTime:  info.ModTime().UnixMilli(),
		IsDir:       info.IsDir(),
		MD5:         md5,
		Thumbnail:   thumbnail,
		IsCollected: isCollected,
	}
}

func GetImageThumbnail(imagePath string, md5 string) (string, error) {
	thumbnail := filepath.Join(filepath.Dir(imagePath), ".thumbnail", fmt.Sprintf("%s.JPEG", md5))
	if _, err := os.Stat(thumbnail); err == nil {
		return thumbnail, nil
	} else {
		return "", err
	}
}

func GetFileInfo(path string) (File, error) {
	info, err := os.Stat(path)
	if err != nil {
		return File{}, err
	}
	return PackFileByInfo(path, info, GetDeviceID()), nil
}

func GetFileInfoWithNID(path string, nid string) (File, error) {
	info, err := os.Stat(path)
	if err != nil {
		return File{}, err
	}
	return PackFileByInfo(path, info, nid), nil
}

// CreateFile creates a new file at the specified path
func CreateFile(path string) (File, error) {
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		return PackFileByInfo(path, info, GetDeviceID()), nil
	}
	dirPath := filepath.Dir(path)
	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		os.MkdirAll(dirPath, 0777)
	}
	_, err = os.Create(path)
	if err != nil {
		return File{}, err
	}
	info, err = os.Stat(path)
	if err != nil {
		return File{}, err
	}
	return PackFileByInfo(path, info, GetDeviceID()), nil
}

func CreateDirectory(path string) (File, error) {
	return CreateDirectoryPerm(path, os.ModePerm)
}

func CreateDirectoryPerm(path string, perm os.FileMode) (File, error){
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return PackFileByInfo(path, info, GetDeviceID()), nil
	}
	err = os.MkdirAll(path, perm)
	if err != nil {
		return File{}, err
	}
	info, err = os.Stat(path)
	if err != nil {
		return File{}, err
	}
	return PackFileByInfo(path, info, GetDeviceID()), nil
}

// DeleteFile deletes the file at the specified path
func DeleteFile(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	// return os.RemoveAll(path)
	AddRecentDeleteFiles([]string{path})
	return nil
}

func MoveFile(oldPath string, newPath string) (File, error) {
	if err := os.Rename(oldPath, newPath); err != nil {
		return File{}, err
	}
	fi, err := os.Stat(newPath)
	if err != nil {
		return File{}, err
	} else {
		MoveCollect(oldPath, newPath)
		MoveRecentOpenFile(oldPath, newPath)
		MoveRecentAddFile(oldPath, newPath)
		return PackFileByInfo(newPath, fi, GetDeviceID()), nil
	}
}

func BatchMoveFiles(oldPaths []string, newDirPath string) error {
	if _, err := os.Stat(newDirPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(newDirPath, os.ModePerm)
		} else {
			return err
		}
	}
	for _, oldPath := range oldPaths {
		newPath := filepath.Join(newDirPath, filepath.Base(oldPath))
		if os.Rename(oldPath, newPath) == nil {
			MoveCollect(oldPath, newPath)
			MoveRecentOpenFile(oldPath, newPath)
			MoveRecentAddFile(oldPath, newPath)
		}
	}
	return nil
}

func GetFileMD5(path string) (string, error) {
	nBytes := 20
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	var bytes []byte
	if fileInfo.Size() <= int64(nBytes)*2 {
		bytes = make([]byte, fileInfo.Size())
		if _, err = file.Read(bytes); err != nil {
			return "", err
		}
	} else {
		firstBytes := make([]byte, nBytes)
		if _, err = file.Read(firstBytes); err != nil {
			return "", err
		}
		if _, err = file.Seek(-int64(nBytes), io.SeekEnd); err != nil {
			return "", err
		}
		lastBytes := make([]byte, nBytes)
		if _, err = file.Read(lastBytes); err != nil {
			return "", err
		}
		bytes = append(firstBytes, lastBytes...)
	}
	hash := md5.Sum(bytes)
	md5 := hex.EncodeToString(hash[:])
	return md5, nil
}
