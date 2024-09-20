package models

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/ulangch/nas_desktop_app/backend/macro"
)

const DEFAULT_DEVICE_NAME = "我的私有云"

type Device struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Disks []File `json:"disks"`
	PSID  string `json:"ps_id"`
	PSDir File   `json:"ps_dir"`
}

type Disk struct {
	ID   string `json:"id"`
	Name string `json:"name"` // Disk name
	Path string `json:"path"` // Absolute disk dir path
}

func InitializeDeviceID() {
	id, err := GetKV(KV_KEY_DEVICE_ID)
	if err != nil {
		log.Printf("InitializeDeviceID GetKV failed: %s", err.Error())
	}
	if id == "" {
		PutKV(KV_KEY_DEVICE_ID, uuid.NewString())
	}
}

func GetDeviceID() string {
	id, _ := GetKV(KV_KEY_DEVICE_ID)
	return id
}

func GetDeviceInfo() (Device, error) {
	nid, err := GetKV(KV_KEY_DEVICE_ID)
	if err != nil {
		return Device{}, errors.New("no device")
	}
	name, _ := GetKV(KV_KEY_DEVICE_NAME)
	if name == "" {
		name = DEFAULT_DEVICE_NAME
	}
	disks, err := GetDeviceDisks()
	if err != nil {
		return Device{}, errors.New("no disk")
	}
	var diskFiles []File
	for _, disk := range disks {
		realPath := disk.Path
		info, err := os.Stat(realPath)
		if err != nil && os.IsNotExist(err) {
			os.MkdirAll(realPath, os.ModePerm)
			info, err = os.Stat(realPath)
		}
		if err != nil || !info.IsDir() {
			log.Printf("disk invalid: %s", realPath)
			continue
		}
		file, err := PackFile4(realPath, info, nid, disk.ID)
		if err != nil {
			log.Printf("pack failed: %s", err.Error())
			continue
		}
		usage, _ := macro.GetDiskUsage(realPath)
		file.FreeVolume = usage.Free
		file.TotalVolume = usage.Total
		diskFiles = append(diskFiles, file)
	}
	psId := GetPrivateSpaceSid()
	psDir, _ := GetPrivateSpaceDir()
	return Device{ID: nid, Name: name, Disks: diskFiles, PSID: psId, PSDir: psDir}, nil
}

func UpdateDeviceName(name string) error {
	return PutKV(KV_KEY_DEVICE_NAME, name)
}

func GetDeviceDisks() ([]Disk, error) {
	var disks []Disk
	if disksJson, _ := GetKV(KV_KEY_DEVICE_DISKS); disksJson != "" {
		if json.Unmarshal([]byte(disksJson), &disks) != nil {
			return nil, errors.New("unmarshal failed")
		}
	}
	if len(disks) <= 0 {
		return nil, errors.New("no disk")
	}
	return disks, nil
}

func GetDeviceDiskId(realPath string) (string, error) {
	disks, err := GetDeviceDisks()
	if err != nil {
		return "", errors.New("no disk")
	}
	for _, disk := range disks {
		if strings.HasPrefix(realPath, disk.Path) {
			return disk.ID, nil
		}
	}
	return "", errors.New("disk not found")
}

// Absolute path for [id]
// windows: F:\Storage\云空间
// darwin: /Users/ulangch/云空间
func GetDeviceDiskPath(id string) (string, error) {
	return GetKV(KV_KEY_DEVICE_DISK_PREFIX + id)
}

func AddDiskPath(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, 0777)
		info, err = os.Stat(path)
	}
	if err != nil {
		return errors.New("disk invalid")
	}
	if !info.IsDir() {
		return errors.New("disk not a directory")
	}
	var disks []Disk
	if disksJson, _ := GetKV(KV_KEY_DEVICE_DISKS); disksJson != "" {
		if json.Unmarshal([]byte(disksJson), &disks) != nil {
			return errors.New("unmarshal failed")
		}
	}
	diskId := GetStringMD5(path)
	for _, disk := range disks {
		if disk.ID == diskId || disk.Path == path {
			return errors.New("disk already added")
		}
	}
	disks = append(disks, Disk{ID: diskId, Name: filepath.Base(path), Path: path})
	disksJsonBytes, err := json.Marshal(disks)
	if err != nil {
		return errors.New("marshal failed")
	}
	if PutKV(KV_KEY_DEVICE_DISK_PREFIX+diskId, path) != nil {
		return errors.New("store disk failed")
	}
	if PutKV(KV_KEY_DEVICE_DISKS, string(disksJsonBytes)) != nil {
		return errors.New("store disk failed")
	} else {
		return nil
	}
}

func UpdateDiskPath(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, 0777)
		info, err = os.Stat(path)
	}
	if err != nil {
		return errors.New("disk invalid")
	}
	if !info.IsDir() {
		return errors.New("disk not a directory")
	}
	var disks []Disk
	if disksJson, _ := GetKV(KV_KEY_DEVICE_DISKS); disksJson != "" {
		if json.Unmarshal([]byte(disksJson), &disks) != nil {
			return errors.New("unmarshal failed")
		}
	}
	diskId := GetStringMD5(path)
	for _, disk := range disks {
		if disk.ID == diskId || disk.Path == path {
			return errors.New("disk already added")
		}
	}
	newDisks := []Disk{
		{ID: diskId, Name: filepath.Base(path), Path: path},
	}
	disksJsonBytes, err := json.Marshal(newDisks)
	if err != nil {
		return errors.New("marshal failed")
	}
	if PutKV(KV_KEY_DEVICE_DISK_PREFIX+diskId, path) != nil {
		return errors.New("store disk failed")
	}
	if PutKV(KV_KEY_DEVICE_DISKS, string(disksJsonBytes)) != nil {
		return errors.New("store disk failed")
	}
	// 清除记录
	for _, disk := range disks {
		DeleteKV(KV_KEY_DEVICE_DISK_PREFIX + disk.ID)
	}
	return nil
}

func RemoveDiskPath(path string) error {
	var disks []Disk
	if disksJson, _ := GetKV(KV_KEY_DEVICE_DISKS); disksJson != "" {
		if json.Unmarshal([]byte(disksJson), &disks) != nil {
			return errors.New("unmarshal failed")
		}
	}
	var existDisk Disk
	var newDisks []Disk
	for _, disk := range disks {
		if disk.Path == path {
			existDisk = disk
		} else {
			newDisks = append(newDisks, disk)
		}
	}
	disksJsonBytes, err := json.Marshal(newDisks)
	if err != nil {
		return errors.New("marshal failed")
	}
	if PutKV(KV_KEY_DEVICE_DISKS, string(disksJsonBytes)) != nil {
		return errors.New("store disk failed")
	} else {
		DeleteKV(KV_KEY_DEVICE_DISK_PREFIX + existDisk.ID)
		return nil
	}
}
