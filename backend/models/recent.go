package models

import (
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type RecentOpenFile struct {
	Path string `gorm:"primaryKey"`
	Time int64
}

type RecentAddFile struct {
	Path string `gorm:"primaryKey"`
	Time int64
}

type RecentStore struct {
	db *gorm.DB
}

var recentStore *RecentStore

func InitializeRecentDB(db *gorm.DB) {
	recentStore = &RecentStore{db: db}
}

func AddRecentOpenFiles(paths []string) {
	for index, path := range paths {
		open := RecentOpenFile{Path: path, Time: time.Now().UnixMilli() + int64(index)}
		recentStore.db.Save(&open)
	}
}

func ListRecentOpenFiles() ([]File, error) {
	var opens []RecentOpenFile
	if err := recentStore.db.Find(&opens).Error; err != nil {
		return nil, err
	}
	var files []File
	nid := GetDeviceID()
	for _, open := range opens {
		if file, err := GetFileInfoWithNID(open.Path, nid); err == nil {
			file.UpdateTime = open.Time
			files = append(files, file)
		}
	}
	return files, nil
}

func AddRecentAddFiles(paths []string) {
	for index, path := range paths {
		add := RecentAddFile{Path: path, Time: time.Now().UnixMilli() + int64(index)}
		recentStore.db.Save(&add)
	}
}

func ListRecentAddFiles() ([]File, error) {
	var adds []RecentAddFile
	if err := recentStore.db.Find(&adds).Error; err != nil {
		return nil, err
	}
	var files []File
	nid := GetDeviceID()
	for _, add := range adds {
		if file, err := GetFileInfoWithNID(add.Path, nid); err == nil {
			file.UpdateTime = add.Time
			files = append(files, file)
		}
	}
	return files, nil
}

func AddRecentDeleteFiles(paths []string) {
	deleteDirPath, _ := GetRecentDeleteDirPath()
	for _, path := range paths {
		fi, err := os.Stat(path)
		if err != nil {
			continue
		}
		dir := filepath.Dir(path)
		if deleteDirPath == "" || dir == deleteDirPath {
			os.RemoveAll(path)
		} else {
			deletePath := filepath.Join(deleteDirPath, fi.Name())
			os.Rename(path, deletePath)
		}
	}
}

func ListRecentDeleteFiles() ([]File, error) {
	deleteDirPath, err := GetRecentDeleteDirPath()
	if err != nil {
		return nil, err
	}
	return ListFiles(deleteDirPath)
}

func GetRecentDeleteDirPath() (string, error) {
	disks, err := GetDeviceDisks()
	if err != nil {
		return "", err
	}
	dirPath := filepath.Join(disks[0].Path, ".delete")
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModePerm)
	}
	return dirPath, nil
}
