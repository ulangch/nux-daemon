package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ulangch/nas_desktop_app/backend/macro"
	"github.com/ulangch/nas_desktop_app/backend/utils"
)

var bootTime int64 = 0

func Initialize() {
	id, err := utils.GetKV(utils.KV_KEY_NUX_DEVICE_ID)
	if err != nil {
		log.Printf("core_service initialize GetKV failed: %s", err.Error())
	}
	if id == "" {
		utils.PutKV(utils.KV_KEY_NUX_DEVICE_ID, uuid.NewString())
	}
	bootTime = time.Now().UnixMilli()
}

func GetDeviceId() (string, error) {
	id, err := utils.GetKV(utils.KV_KEY_NUX_DEVICE_ID)
	if err != nil {
		log.Printf("core_service GetDeviceId GetKV failed: %s", err.Error())
		return "", err
	}
	if id == "" {
		id = uuid.NewString()
		utils.PutKV(utils.KV_KEY_NUX_DEVICE_ID, id)
	}
	return id, nil
}

func GetDevice() (NuxDevice, error) {
	id, err := GetDeviceId()
	if err != nil {
		return NuxDevice{}, errors.New("core_service GetDevice, GetDeviceId failed")
	}
	name, _ := utils.GetKV(utils.KV_KEY_NUX_DEVICE_NAME)
	if name == "" {
		name = utils.DEFAULT_DEVICE_NAME
	}
	disks, err := GetDisks()
	if err != nil {
		log.Printf("core_service GetDevice, no disk")
	}
	// psId := GetPrivateSpaceSid()
	// psDir, _ := GetPrivateSpaceDir()
	var diskTotalBytes int64
	var diskFreeBytes int64
	var iterates []NuxDisk
	for _, disk := range disks {
		var counted = false
		diskPlatformPath := utils.ToPlatformPath(disk.UnixAbsolute)
		for _, iterated := range iterates {
			iteratedPlatformPath := utils.ToPlatformPath(iterated.UnixAbsolute)
			if macro.IsSameVolume(diskPlatformPath, iteratedPlatformPath) {
				counted = true
				break
			}
		}
		if !counted {
			diskTotalBytes += disk.DiskTotalBytes
			diskFreeBytes += disk.DiskFreeBytes
		}
		iterates = append(iterates, disk)
	}
	ipv4Addr, _ := utils.GetLocalIPv4()
	memoryTotalBytes, memoryFreeBytes := macro.GetSystemMemory()
	nuxDevice := NuxDevice{
		ID:               id,
		Name:             name,
		BootTime:         bootTime,
		Address:          ipv4Addr,
		Disks:            disks,
		DiskTotalBytes:   diskTotalBytes,
		DiskFreeBytes:    diskFreeBytes,
		MemoryTotalBytes: memoryTotalBytes,
		MemoryFreeBytes:  memoryFreeBytes,
		CPURate:          int32(macro.GetSystemCpuRate()),
		CPUTemp:          int32(macro.GetSystemCpuTemperature()),
	}
	return nuxDevice, nil
}

func GetDisks() ([]NuxDisk, error) {
	var diskInStore []NuxDisk
	if disksJson, _ := utils.GetKV(utils.KV_KEY_NUX_DEVICE_DISKS); disksJson != "" {
		if json.Unmarshal([]byte(disksJson), &diskInStore) != nil {
			return nil, errors.New("core_service GetDisks, unmarshal failed")
		}
	}
	if len(diskInStore) <= 0 {
		return nil, errors.New("core_service GetDisks, no disk")
	}
	var disks []NuxDisk
	for _, disk := range diskInStore {
		platformPath := utils.ToPlatformPath(disk.UnixAbsolute)
		info, err := os.Stat(platformPath)
		if err != nil && os.IsNotExist(err) {
			os.MkdirAll(platformPath, os.ModePerm)
			info, err = os.Stat(platformPath)
		}
		if err != nil || !info.IsDir() {
			log.Printf("core_service GetDevice, disk invalid: %s", disk.UnixAbsolute)
			continue
		}
		usage, _ := macro.GetDiskUsage(platformPath)
		disk.DiskTotalBytes = usage.Total
		disk.DiskFreeBytes = usage.Free
		disks = append(disks, disk)
	}
	return disks, nil
}

func GetDiskById(diskId string) (NuxDisk, error) {
	disks, err := GetDisks()
	if err != nil {
		return NuxDisk{}, err
	}
	for _, disk := range disks {
		if disk.ID == diskId {
			return disk, nil
		}
	}
	return NuxDisk{}, fmt.Errorf("core_device GetDisk, no disk, diskId=%s", diskId)
}

func GetDiskByAbsolutePath(absolutePath string) (NuxDisk, error) {
	unixAbsolute := utils.ToUnixPath(absolutePath)
	disks, err := GetDisks()
	if err != nil {
		return NuxDisk{}, err
	}
	for _, disk := range disks {
		if strings.Contains(unixAbsolute, disk.UnixAbsolute) {
			return disk, nil
		}
	}
	return NuxDisk{}, fmt.Errorf("core_device GetDisk, no disk, absolutePath=%s", absolutePath)
}

func AddDiskByPlatformPath(platformPath string) error {
	info, err := os.Stat(platformPath)
	if os.IsNotExist(err) {
		os.MkdirAll(platformPath, 0777)
		info, err = os.Stat(platformPath)
	}
	if err != nil {
		return errors.New("core_service AddDisk, path invalid")
	}
	if !info.IsDir() {
		return errors.New("core_service AddDisk, path not a directory")
	}
	var disks []NuxDisk
	if disksJson, _ := utils.GetKV(utils.KV_KEY_NUX_DEVICE_DISKS); disksJson != "" {
		if json.Unmarshal([]byte(disksJson), &disks) != nil {
			return errors.New("core_service AddDisk, unmarshal failed")
		}
	}
	unixAbsolute := utils.ToUnixPath(platformPath)
	diskId := utils.GetStringMD5(unixAbsolute)
	for _, disk := range disks {
		if disk.ID == diskId || disk.UnixAbsolute == unixAbsolute {
			return errors.New("core_service AddDisk, disk already added")
		}
	}
	disk := NuxDisk{ID: diskId, Name: info.Name(), UnixPath: "/" + diskId, UnixAbsolute: unixAbsolute}
	disks = append(disks, disk)
	disksJsonBytes, err := json.Marshal(disks)
	if err != nil {
		return errors.New("core_service AddDisk, marshal failed")
	}
	if utils.PutKV(utils.KV_KEY_NUX_DEVICE_DISK_PREFIX+diskId, unixAbsolute) != nil {
		return errors.New("core_service AddDisk, store disk path failed")
	}
	if utils.PutKV(utils.KV_KEY_NUX_DEVICE_DISKS, string(disksJsonBytes)) != nil {
		return errors.New("core_service AddDisk, store disk failed")
	} else {
		return nil
	}
}

func RemoveDiskByPlatformPath(platformPath string) error {
	var disks []NuxDisk
	if disksJson, _ := utils.GetKV(utils.KV_KEY_NUX_DEVICE_DISKS); disksJson != "" {
		if json.Unmarshal([]byte(disksJson), &disks) != nil {
			return errors.New("core_service RemoveDisk, unmarshal failed")
		}
	}
	unixAbsolute := utils.ToUnixPath(platformPath)
	var existDisk NuxDisk
	var newDisks []NuxDisk
	for _, disk := range disks {
		if disk.UnixAbsolute == unixAbsolute {
			existDisk = disk
		} else {
			newDisks = append(newDisks, disk)
		}
	}
	disksJsonBytes, err := json.Marshal(newDisks)
	if err != nil {
		return errors.New("core_service RemoveDisk, marshal failed")
	}
	if utils.PutKV(utils.KV_KEY_NUX_DEVICE_DISKS, string(disksJsonBytes)) != nil {
		return errors.New("core_service RemoveDisk, store disk failed")
	} else {
		utils.DeleteKV(utils.KV_KEY_NUX_DEVICE_DISK_PREFIX + existDisk.ID)
		return nil
	}
}
