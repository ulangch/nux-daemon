package models

import (
	"errors"
	"os"
	"path/filepath"
)

func UpdateGalleryDir(path string) (File, error) {
	fi, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
		fi, err = os.Stat(path)
	}
	if err != nil {
		return File{}, err
	} else if err := PutKV(KV_KEY_GALLERY_DIR, path); err != nil {
		return File{}, err
	} else {
		return PackFile2(path, fi)
	}
}

func GetGalleryDir(clientModel string) (File, error) {
	var path string
	var err error
	if path, err = GetKV(KV_KEY_GALLERY_DIR); err != nil || path == "" {
		var device Device
		if device, err = GetDeviceInfo(); err == nil && len(device.Disks) > 0 {
			disk := device.Disks[0]
			path = filepath.Join(disk.Path, clientModel, "云相册")
		}
	}
	if err != nil {
		return File{}, err
	} else if path == "" {
		return File{}, errors.New("unknown error")
	} else {
		fi, err := os.Stat(path)
		if err != nil && os.IsNotExist(err) {
			os.MkdirAll(path, os.ModePerm)
			fi, err = os.Stat(path)
		}
		if err != nil {
			return File{}, err
		} else {
			return PackFile2(path, fi)
		}
	}
}

func ListGalleryFiles(clientModel string) ([]File, error) {
	galleryDir, err := GetGalleryDir(clientModel)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(galleryDir.Path)
	if err != nil {
		return nil, err
	}
	nid := GetDeviceID()
	diskId, err := GetDeviceDiskId(galleryDir.Path)
	if err != nil {
		return nil, err
	}
	var imageFiles []File
	for _, entry := range entries {
		if !FilterFile(entry.Name()) {
			continue
		}
		if entry.IsDir() {
			// Album
			albumPath := filepath.Join(galleryDir.Path, entry.Name())
			if albumEntries, err := os.ReadDir(albumPath); err == nil {
				for _, albumEntry := range albumEntries {
					if !albumEntry.IsDir() && (IsImage(albumEntry.Name()) || IsVideo(albumEntry.Name())) {
						if info, err := albumEntry.Info(); err == nil {
							imagePath := filepath.Join(albumPath, info.Name())
							if imageFile, err := PackFile4(imagePath, info, nid, diskId); err == nil {
								imageFiles = append(imageFiles, imageFile)
							}
						}
					}
				}
			}
		} else if IsImage(entry.Name()) || IsVideo(entry.Name()) {
			// Image
			if info, err := entry.Info(); err == nil {
				imagePath := filepath.Join(galleryDir.Path, info.Name())
				if imageFile, err := PackFile4(imagePath, info, nid, diskId); err == nil {
					imageFiles = append(imageFiles, imageFile)
				}
			}
		}
	}
	return imageFiles, nil
}
