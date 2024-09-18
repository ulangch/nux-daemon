package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ulangch/nas_desktop_app/backend/macro"
)

type File struct {
	Nid         string `json:"nid"`
	Name        string `json:"name"`
	Path        string `json:"-"`   // real path
	Url         string `json:"url"` // client url
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

// clientPath: [diskId]/[path]
func GetRealPath(clientPath string) (string, error) {
	// clientPath, _ = strings.CutPrefix(clientPath, "/")
	segments := strings.Split(clientPath, "/")
	if len(segments) <= 0 {
		return "", errors.New("path invalid")
	}
	diskId := segments[0]
	// windows: F:\Storage\云空间
	// darwin: /Users/ulangch/云空间
	diskPath, err := GetDeviceDiskPath(diskId)
	if err != nil {
		return "", errors.New("disk not found")
	}
	relativePath, _ := strings.CutPrefix(clientPath, diskId)
	realPath := diskPath + macro.DecodeFilePath(relativePath)
	return realPath, nil
}

func GetRealPaths(clientPaths []string) ([]string, error) {
	var realPaths []string
	for _, clientPath := range clientPaths {
		if realPath, err := GetRealPath(clientPath); err != nil {
			return nil, fmt.Errorf("path invalid, path=%s", clientPath)
		} else {
			realPaths = append(realPaths, realPath)
		}
	}
	return realPaths, nil
}

func GetClientPath(realPath string) (string, error) {
	disks, err := GetDeviceDisks()
	if err != nil {
		return "", errors.New("path invalid")
	}
	for _, disk := range disks {
		if strings.HasPrefix(realPath, disk.Path) {
			relativePath, _ := strings.CutPrefix(realPath, disk.Path)
			return disk.ID + macro.EncodeFilePath(relativePath), nil
		}
	}
	return "", errors.New("disk not found")
}

func GetClientPath2(realPath string, diskId string) (string, error) {
	diskPath, err := GetDeviceDiskPath(diskId)
	if err != nil {
		return "", errors.New("disk not found")
	}
	if relativePath, hasFound := strings.CutPrefix(realPath, diskPath); !hasFound {
		return "", errors.New("path invalid")
	} else {
		return diskId + macro.EncodeFilePath(relativePath), nil
	}
}

func GetClientPaths(realPaths []string) ([]string, error) {
	var clientPaths []string
	for _, realPath := range realPaths {
		if clientPath, err := GetClientPath(realPath); err != nil {
			return nil, fmt.Errorf("path invalid, path=%s", realPath)
		} else {
			clientPaths = append(clientPaths, clientPath)
		}
	}
	return clientPaths, nil
}

func GetClientPathsWithDiskId(realPaths []string, diskId string) ([]string, error) {
	diskPath, err := GetDeviceDiskPath(diskId)
	if err != nil {
		return nil, errors.New("disk not found")
	}
	var clientPaths []string
	for _, realPath := range realPaths {
		if relativePath, hasFound := strings.CutPrefix(realPath, diskPath); !hasFound {
			return nil, fmt.Errorf("path invalid, path=%s", realPath)
		} else {
			clientPath := diskId + macro.EncodeFilePath(relativePath)
			clientPaths = append(clientPaths, clientPath)
		}
	}
	return clientPaths, nil
}

func GetClientUrl(realPath string, nid string) (string, error) {
	if clientPath, err := GetClientPath(realPath); err != nil {
		return "", err
	} else {
		return "nas://" + nid + "/" + clientPath, nil
	}
}

func GetClientUrlWithDiskId(realPath string, nid string, diskId string) (string, error) {
	if clientPath, err := GetClientPath2(realPath, diskId); err != nil {
		return "", err
	} else {
		return "nas://" + nid + "/" + clientPath, nil
	}
}

// ListFiles lists all files in the specified directory
// path: real path
func ListFiles(path string) ([]File, error) {
	var files []File
	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	nid := GetDeviceID()
	diskId, err := GetDeviceDiskId(path)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if !FilterFile(entry.Name()) {
			continue
		}
		file, err := PackFile4(filepath.Join(path, info.Name()), info, nid, diskId)
		if err != nil {
			log.Printf("pack failed: %s", err.Error())
			continue
		}
		files = append(files, file)
	}
	return files, nil
}

// path: real path
func ListTypeFiles(path string, filterType string, maxDepth int) ([]File, error) {
	var files []File
	if err := ListTypeFilesRecursive(&files, path, filterType, 0, maxDepth); err != nil {
		return nil, err
	} else {
		return files, nil
	}
}

// List [filterType] files in [dir], depth start with 0, maxDepth should be equal or bigger than 0
// path: real path
func ListTypeFilesRecursive(result *[]File, path string, filterType string, depth int, maxDepth int) error {
	if depth > maxDepth {
		return nil
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}
	nid := GetDeviceID()
	diskId, err := GetDeviceDiskId(path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}
		path := filepath.Join(path, entry.Name())
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
		file, err := PackFile4(path, info, nid, diskId)
		if err != nil {
			log.Printf("pack failed: %s", err.Error())
			continue
		}
		*result = append(*result, file)
	}
	return nil
}

func PackFile1(path string) {

}

func PackFile2(path string, info os.FileInfo) (File, error) {
	nid := GetDeviceID()
	return PackFile3(path, info, nid)
}

func PackFile3(path string, info os.FileInfo, nid string) (File, error) {
	diskId, err := GetDeviceDiskId(path)
	if err != nil {
		return File{}, err
	}
	return PackFile4(path, info, nid, diskId)
}

func PackFile4(path string, info os.FileInfo, nid string, diskId string) (File, error) {
	var md5 string
	if !info.IsDir() {
		md5, _ = GetFileMD5(path)
	}
	isCollected := IsCollected(path)
	clientUrl, err := GetClientUrlWithDiskId(path, nid, diskId)
	if err != nil {
		return File{}, err
	}
	return File{
		Nid:         nid,
		Name:        info.Name(),
		Path:        path,
		Url:         clientUrl,
		Size:        info.Size(),
		UpdateTime:  info.ModTime().UnixMilli(),
		IsDir:       info.IsDir(),
		MD5:         md5,
		IsCollected: isCollected,
	}, nil
}

// imagePath: real path
func GetImageThumbnail(imagePath string, md5 string) (string, error) {
	thumbnail := filepath.Join(filepath.Dir(imagePath), ".thumbnail", fmt.Sprintf("%s.JPEG", md5))
	if _, err := os.Stat(thumbnail); err == nil {
		return thumbnail, nil
	} else {
		return "", err
	}
}

// path: real path
func GetFileInfo(path string) (File, error) {
	info, err := os.Stat(path)
	if err != nil {
		return File{}, err
	}
	return PackFile2(path, info)
}

// path: real path
func GetFileInfo2(path string, nid string) (File, error) {
	info, err := os.Stat(path)
	if err != nil {
		return File{}, err
	}
	return PackFile3(path, info, nid)
}

// CreateFile creates a new file at the specified path
// path: real path
func CreateFile(path string) (File, error) {
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		return PackFile2(path, info)
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
	return PackFile2(path, info)
}

// path: real path
func CreateDirectory(path string) (File, error) {
	return CreateDirectoryPerm(path, os.ModePerm)
}

// path: real path
func CreateDirectoryPerm(path string, perm os.FileMode) (File, error) {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return PackFile2(path, info)
	}
	err = os.MkdirAll(path, perm)
	if err != nil {
		return File{}, err
	}
	info, err = os.Stat(path)
	if err != nil {
		return File{}, err
	}
	return PackFile2(path, info)
}

// DeleteFile deletes the file at the specified path
// path: real path
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

// path: real path
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
		return PackFile2(newPath, fi)
	}
}

// path: real path
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

// path: real path
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
