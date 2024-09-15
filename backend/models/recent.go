package models

import (
	"fmt"
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

type RecentDeleteFile struct {
	Path         string `gorm:"primaryKey"`
	OriginalPath string
	Time         int64
}

type RecentStore struct {
	db *gorm.DB
}

var recentStore *RecentStore

func InitializeRecentDB(db *gorm.DB) {
	recentStore = &RecentStore{db: db}
}

func AddRecentOpenFiles(paths []string) {
	psDirPath := GetPrivateSpaceDirPath()
	for index, path := range paths {
		if !IsInPrivateSpaceDir2(path, psDirPath) {
			open := RecentOpenFile{Path: path, Time: time.Now().UnixMilli() + int64(index)}
			recentStore.db.Save(&open)
		}
	}
}

func MoveRecentOpenFile(oldPath string, newPath string) {
	var recentOpen RecentOpenFile
	if recentStore.db.First(&recentOpen, "path = ?", oldPath).Error == nil {
		recentStore.db.Delete(&RecentOpenFile{}, "path = ?", oldPath)
		if !IsInPrivateSpaceDir(newPath) {
			recentOpen = RecentOpenFile{Path: newPath, Time: recentOpen.Time}
			recentStore.db.Save(&recentOpen)
		}
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
	psDirPath := GetPrivateSpaceDirPath()
	for index, path := range paths {
		if !IsInPrivateSpaceDir2(path, psDirPath) {
			add := RecentAddFile{Path: path, Time: time.Now().UnixMilli() + int64(index)}
			recentStore.db.Save(&add)
		}
	}
}

func MoveRecentAddFile(oldPath string, newPath string) {
	var recentAdd RecentAddFile
	if recentStore.db.First(&recentAdd, "path = ?", oldPath).Error == nil {
		recentStore.db.Delete(&RecentAddFile{}, "path = ?", oldPath)
		if !IsInPrivateSpaceDir(newPath) {
			recentAdd = RecentAddFile{Path: newPath, Time: recentAdd.Time}
			recentStore.db.Save(&recentAdd)
		}
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
			recentStore.db.Delete(&RecentDeleteFile{}, "path = ?", path)
		} else {
			deletePath := filepath.Join(deleteDirPath, fi.Name())
			os.Rename(path, deletePath)
			recentDelete := RecentDeleteFile{Path: deletePath, OriginalPath: path, Time: time.Now().UnixMilli()}
			recentStore.db.Save(&recentDelete)
		}
	}
}

func RecoverRecentDeleteFiles(paths []string) error {
	var recentDeletes []RecentDeleteFile
	if err := recentStore.db.Find(&recentDeletes, "path IN (?)", paths).Error; err != nil {
		return err
	}
	for _, delete := range recentDeletes {
		os.Rename(delete.Path, delete.OriginalPath)
	}
	recentStore.db.Delete(&RecentDeleteFile{}, "path IN (?)", paths)
	return nil
}

func ListRecentDeleteFiles() ([]File, error) {
	var recentDeletes []RecentDeleteFile
	if err := recentStore.db.Find(&recentDeletes).Error; err != nil {
		return nil, err
	}
	nid := GetDeviceID()
	var files []File
	for _, delete := range recentDeletes {
		if file, err := GetFileInfoWithNID(delete.Path, nid); err == nil {
			file.UpdateTime = delete.Time
			file.GhostUrl = fmt.Sprintf("nas://%s%s", nid, delete.OriginalPath)
			files = append(files, file)
		}
	}
	return files, nil
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
