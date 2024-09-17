package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/ulangch/nas_desktop_app/backend/macro"
)

type Device struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Disks []File `json:"disks"`
	PSID  string `json:"ps_id"`
	PSDir File   `json:"ps_dir"`
}

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

func GetDeviceID() string {
	id, _ := GetKV(KV_KEY_ID)
	return id
}

func GetDeviceDisks() ([]File, error) {
	pathsJson, _ := GetKV(KV_KEY_DISKS)
	var paths []string
	if pathsJson != "" {
		err := json.Unmarshal([]byte(pathsJson), &paths)
		if err != nil {
			log.Printf("GetDeviceInfo Unmarshal failed: %s", err.Error())
			return nil, err
		}
	}
	nid := GetDeviceID()
	var disks []File
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil && os.IsNotExist(err) {
			os.MkdirAll(path, os.ModePerm)
			info, err = os.Stat(path)
		}
		if err != nil {
			log.Printf("GetDeviceInfo stat failed: %s", err.Error())
			continue
		}
		if !info.IsDir() {
			log.Printf("GetDeviceInfo not a directory: %s", path)
			continue
		}

		// Calculate total and free space
		// var stat syscall.Statfs_t
		// var total uint64
		// var free uint64
		// err = syscall.Statfs(path, &stat)
		// if err == nil {
		// 	total = stat.Blocks * uint64(stat.Bsize)
		// 	free = stat.Bfree * uint64(stat.Bsize)
		// } else {
		// 	log.Printf("GetDeviceInfo get volume failed: %s", err.Error())
		// }

		path = macro.EncodeFilePath(path)

		diskUsage, _ := macro.GetDiskUsage(path)
		disks = append(disks, File{
			Nid:         nid,
			Name:        info.Name(),
			Path:        path,
			Url:         fmt.Sprintf("nas://%s%s", nid, path),
			Size:        info.Size(),
			UpdateTime:  info.ModTime().UnixMilli(),
			IsDir:       true,
			FreeVolume:  int64(diskUsage.Free),
			TotalVolume: int64(diskUsage.Total),
		})
	}
	if len(disks) <= 0 {
		return nil, errors.New("Device disks empty")
	} else {
		return disks, nil
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
	disks, err := GetDeviceDisks()
	if err != nil {
		log.Printf("GetDeviceInfo GetDisks failed: %s", err.Error())
		return Device{}, err
	}
	psId := GetPrivateSpaceSid()
	psDir, _ := GetPrivateSpaceDir()
	return Device{ID: id, Name: name, Disks: disks, PSID: psId, PSDir: psDir}, nil
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
