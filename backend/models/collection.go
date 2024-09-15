package models

import (
	"time"

	"gorm.io/gorm"
)

type Collection struct {
	Path string `gorm:"primaryKey"`
	Time int64
}

type CollectionStore struct {
	db *gorm.DB
}

var colStore *CollectionStore

func InitializeColStore(db *gorm.DB) {
	colStore = &CollectionStore{db: db}
}

func CollectFiles(paths []string) {
	psDirPath := GetPrivateSpaceDirPath()
	for index, path := range paths {
		if !IsInPrivateSpaceDir2(path, psDirPath) {
			collection := Collection{Path: path, Time: time.Now().UnixMilli() + int64(index)}
			colStore.db.Save(&collection)
		}
	}
}

func UnCollectFiles(paths []string) {
	for _, path := range paths {
		colStore.db.Delete(&Collection{}, "path = ?", path)
	}
}

func IsCollected(path string) bool {
	var collection Collection
	err := colStore.db.First(&collection, "path = ?", path).Error
	return err == nil
}

func MoveCollect(oldPath string, newPath string) {
	var collection Collection
	if colStore.db.First(&collection, "path = ?", oldPath).Error == nil {
		colStore.db.Delete(&Collection{}, "path = ?", oldPath)
		if !IsInPrivateSpaceDir(newPath) {
			collection = Collection{Path: newPath, Time: collection.Time}
			colStore.db.Save(&collection)
		}
	}
}

func ListCollectFiles() ([]File, error) {
	var collections []Collection
	if err := colStore.db.Find(&collections).Error; err != nil {
		return nil, err
	}
	var files []File
	nid := GetDeviceID()
	for _, collection := range collections {
		if file, err := GetFileInfoWithNID(collection.Path, nid); err == nil {
			file.UpdateTime = collection.Time
			files = append(files, file)
		}
	}
	return files, nil
}
