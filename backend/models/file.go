package models

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
)

type File struct {
	Nid         string `json:"nid"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        int64  `json:"size"`
	UpdateTime  int64  `json:"update_time"`
	IsDir       bool   `json:"is_dir"`
	MD5         string `json:"md5"`
	FreeVolume  int64  `json:"free_volume"`
	TotalVolume int64  `json:"total_volume"`
}

// ListFiles lists all files in the specified directory
func ListFiles(dir string) ([]File, error) {
	var files []File
	entries, err := os.ReadDir(dir)
	log.Printf("ListFiles, dir=%s, entries=%d", dir, len(entries))
	if err != nil {
		return nil, err
	}
	deviceID := GetDeviceID()
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		files = append(files, File{
			Nid:        deviceID,
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

func GetFileInfo(path string, needMD5 bool) (File, error) {
	info, err := os.Stat(path)
	if err != nil {
		return File{}, err
	}
	var md5str string
	if needMD5 && !info.IsDir() {
		md5str, _ = GetFileMD5(path)
	}
	return File{
		Nid:        GetDeviceID(),
		Name:       info.Name(),
		Path:       path,
		Size:       info.Size(),
		UpdateTime: info.ModTime().UnixMilli(),
		IsDir:      info.IsDir(),
		MD5:        md5str,
	}, nil
}

// CreateFile creates a new file at the specified path
func CreateFile(path string) (File, error) {
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		return File{
			Nid:        GetDeviceID(),
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
		Nid:        GetDeviceID(),
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
			Nid:        GetDeviceID(),
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
		Nid:        GetDeviceID(),
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

func MoveFile(oldPath string, newPath string) error {
	return os.Rename(oldPath, newPath)
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
	log.Printf("GetFileMD5, md5=%s, bytes=%x", md5, bytes)
	return md5, nil
}

func WriteFileChunk(path string, offset int64, data []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	log.Printf("WriteFileChunk start, path=%s, offset=%d, data=%d", path, offset, len(data))

	// Move the file pointer to the specified offset
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	// Write the data to the file
	_, err = file.Write(data)
	return err
}
