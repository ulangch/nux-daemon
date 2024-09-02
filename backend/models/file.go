package models

import (
	"os"
	"path/filepath"
)

type File struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        int64  `json:"size"`
	UpdateTime  int64  `json:"update_time"`
	IsDir       bool   `json:"is_dir"`
	FreeVolume  int64  `json:"free_volume"`
	TotalVolume int64  `json:"total_volume"`
}

// ListFiles lists all files in the specified directory
func ListFiles(dir string) ([]File, error) {
	var files []File
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		files = append(files, File{
			Name:       info.Name(),
			Path:       filepath.Join(dir, info.Name()),
			Size:       info.Size(),
			UpdateTime: info.ModTime().UnixMilli(),
			IsDir:      info.IsDir(),
		})
	}
	// 打印日志
	// log.Printf("Files in directory %s: %v\n", dir, files)
	return files, nil
}

func GetFileInfo(path string) (File, error) {
	info, err := os.Stat(path)
	if err != nil {
		return File{}, err
	}
	return File{
		Name:       info.Name(),
		Path:       path,
		Size:       info.Size(),
		UpdateTime: info.ModTime().UnixMilli(),
		IsDir:      info.IsDir(),
	}, nil
}

// CreateFile creates a new file at the specified path
func CreateFile(path string) (File, error) {
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		return File{
			Name:       info.Name(),
			Path:       path,
			Size:       info.Size(),
			UpdateTime: info.ModTime().UnixMilli(),
			IsDir:      info.IsDir(),
		}, nil
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
	return File{
		Name:       info.Name(),
		Path:       path,
		Size:       info.Size(),
		UpdateTime: info.ModTime().UnixMilli(),
		IsDir:      info.IsDir(),
	}, nil
}

func CreateDirectory(path string) (File, error) {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return File{
			Name:       info.Name(),
			Path:       path,
			Size:       info.Size(),
			UpdateTime: info.ModTime().UnixMilli(),
			IsDir:      info.IsDir(),
		}, nil
	}
	err = os.Mkdir(path, 0777)
	if err != nil {
		return File{}, err
	}
	info, err = os.Stat(path)
	if err != nil {
		return File{}, err
	}
	return File{
		Name:       info.Name(),
		Path:       path,
		Size:       info.Size(),
		UpdateTime: info.ModTime().UnixMilli(),
		IsDir:      info.IsDir(),
	}, nil
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
	return os.RemoveAll(path)
}
