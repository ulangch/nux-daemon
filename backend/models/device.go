package models

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"syscall"

	"github.com/google/uuid"
)

type Device struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Disks []File `json:"disks"`
}

const KV_KEY_ID = "KV_KEY_ID"
const KV_KEY_NAME = "KV_KEY_NAME"
const KV_KEY_DISKS = "KV_KEY_DISKS"
const DEFAULT_DEVICE_NAME = "我的私有云"

func InitializeDeviceID() {
	id, err := GetKV(KV_KEY_ID)
	if err != nil {
		log.Printf("InitializeDeviceID GetKV failed: %s", err.Error())
	}
	if id == "" {
		PutKV(KV_KEY_ID, uuid.NewString())
	}
}

func GetDeviceInfo() (Device, error) {
	id, err := GetKV(KV_KEY_ID)
	if err != nil {
		log.Printf("GetDeviceInfo GetID failed: %s", err.Error())
		return Device{}, err
	}
	name, _ := GetKV(KV_KEY_NAME)
	if name == "" {
		name = DEFAULT_DEVICE_NAME
	}
	pathsJson, _ := GetKV(KV_KEY_DISKS)
	var paths []string
	if pathsJson != "" {
		err = json.Unmarshal([]byte(pathsJson), &paths)
		if err != nil {
			log.Printf("GetDeviceInfo Unmarshal failed: %s", err.Error())
			return Device{}, err
		}
	}
	var disks []File
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			log.Printf("GetDeviceInfo stat failed: %s", err.Error())
			continue
		}
		if !info.IsDir() {
			log.Printf("GetDeviceInfo not a directory: %s", path)
			continue
		}

		// Calculate total and free space
		var stat syscall.Statfs_t
		var total uint64
		var free uint64
		err = syscall.Statfs(path, &stat)
		if err == nil {
			total = stat.Blocks * uint64(stat.Bsize)
			free = stat.Bfree * uint64(stat.Bsize)
		} else {
			log.Printf("GetDeviceInfo get volume failed: %s", err.Error())
		}
		disks = append(disks, File{
			Name:        info.Name(),
			Path:        path,
			Size:        info.Size(),
			UpdateTime:  info.ModTime().UnixMilli(),
			IsDir:       true,
			FreeVolume:  int64(free),
			TotalVolume: int64(total),
		})
	}
	return Device{ID: id, Name: name, Disks: disks}, nil
}

func UpdateDeviceName(name string) error {
	return PutKV(KV_KEY_NAME, name)
}

func AddDiskPath(path string, autoCreate bool) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) && autoCreate {
		os.MkdirAll(path, 0777)
		info, err = os.Stat(path)
	}
	if err != nil {
		log.Printf("AddDiskPath Stat failed: %s", err.Error())
		return err
	}
	if !info.IsDir() {
		log.Printf("AddDiskPath path not a directory")
		return errors.New("path not a directory")
	}
	pathsJson, _ := GetKV(KV_KEY_DISKS)
	var paths []string
	if pathsJson != "" {
		err = json.Unmarshal([]byte(pathsJson), &paths)
		if err != nil {
			log.Printf("AddDiskPath Unmarshal failed: %s", err.Error())
			return err
		}
	}
	paths = append(paths, path)
	pathsJsonBytes, err := json.Marshal(paths)
	if err != nil {
		log.Printf("AddDiskPath Marshal failed: %s", err.Error())
		return err
	}
	pathsJson = string(pathsJsonBytes)
	err = PutKV(KV_KEY_DISKS, pathsJson)
	if err != nil {
		log.Printf("AddDiskPath PutKV failed: %s", err.Error())
		return err
	}
	return nil
}

func RemoveDiskPath(path string) error {
	pathsJson, err := GetKV(KV_KEY_DISKS)
	if err != nil {
		log.Printf("RemoveDiskPath GetKV failed: %s", err.Error())
		return err
	}
	var paths []string
	if pathsJson != "" {
		err = json.Unmarshal([]byte(pathsJson), &paths)
	}
	if err != nil {
		log.Printf("RemoveDiskPath Unmarshal failed: %s", err.Error())
		return err
	}
	// paths = paths.remove()
	var newPaths []string
	for _, element := range paths {
		if element != path {
			newPaths = append(newPaths, element)
		}
	}
	pathsJsonBytes, err := json.Marshal(newPaths)
	if err != nil {
		log.Printf("RemoveDiskPath Marshal failed: %s", err.Error())
		return err
	}
	pathsJson = string(pathsJsonBytes)
	err = PutKV(KV_KEY_DISKS, pathsJson)
	if err != nil {
		log.Printf("RemoveDiskPath PutKV failed: %s", err.Error())
		return err
	}
	return nil
}
