package models

import (
	"log"
	"os"

	"gorm.io/gorm"
)

type Collection struct {
	Path string `gorm:"primaryKey"`
}

type CollectionStore struct {
	db *gorm.DB
}

var colStore *CollectionStore

func InitializeColStore(db *gorm.DB) {
	colStore = &CollectionStore{db: db}
}

func CollectFiles(paths []string) {
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			log.Printf("Collect invalid path=%s", path)
			continue
		}
		collection := Collection{Path: path}
		colStore.db.Save(&collection)
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
	UnCollectFiles([]string{oldPath})
	CollectFiles([]string{newPath})
}

func ListCollectFiles() ([]File, error) {
	var paths []string
	if err := colStore.db.Model(&Collection{}).Pluck("Path", &paths).Error; err != nil {
		return nil, err
	}
	var files []File
	nid := GetDeviceID()
	for _, path := range paths {
		if file, err := GetFileInfoWithNID(path, nid); err == nil {
			files = append(files, file)
		}
	}
	return files, nil
}
