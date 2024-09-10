package models

import (
	"encoding/json"
	"errors"
	"log"
	"path/filepath"
	"strings"
)

const BUCKET_VIDEO = "VIDEO"
const BUCKET_AUDIO = "AUDIO"
const BUCKET_DOC = "DOC"
const BUCKET_NOVEL = "NOVEL"

func GetUploadDir(clientModel string) (File, error) {
	dirPath, err := GetKV(KV_KEY_UPLOAD_DIR)
	if err != nil {
		disks, err := GetDeviceDisks()
		if err != nil {
			return File{}, err
		}
		dirPath = filepath.Join(disks[0].Path, clientModel, "云文件")
	}
	return CreateDirectory(dirPath)
}

func SetUploadDir(path string) (File, error) {
	err := PutKV(KV_KEY_UPLOAD_DIR, path)
	if err != nil {
		return File{}, err
	}
	return CreateDirectory(path)
}

func GetBucketDirs(clientModel string, bucket string) ([]File, error) {
	bucketKey := FormatBucketDirsKey(bucket)
	var dirFiles []File
	dirPathsJson, err := GetKV(bucketKey)
	if err != nil || dirPathsJson == "" {
		file, err := GetUploadDir(clientModel)
		if err != nil {
			return nil, err
		} else {
			dirFiles = append(dirFiles, file)
			return dirFiles, nil
		}
	}
	var dirPaths []string
	err = json.Unmarshal([]byte(dirPathsJson), &dirPaths)
	if err != nil {
		return nil, err
	}
	for _, path := range dirPaths {
		if file, err := CreateDirectory(path); err == nil {
			dirFiles = append(dirFiles, file)
		}
	}
	return dirFiles, nil
}

func AddBucketDirs(clientModel string, bucket string, paths []string) ([]File, error) {
	bucketKey := FormatBucketDirsKey(bucket)
	var dirPaths []string
	if dirFiles, err := GetBucketDirs(clientModel, bucket); err == nil {
		for _, file := range dirFiles {
			dirPaths = append(dirPaths, file.Path)
		}
	}
	var addDirFiles []File
	for _, path := range paths {
		if file, err := CreateDirectory(path); err == nil {
			addDirFiles = append(addDirFiles, file)
			dirPaths = append(dirPaths, path)
		} else {
			log.Printf("AddBucketDirs path invalid: %s", path)
		}
	}
	dirPathsJsonBytes, err := json.Marshal(dirPaths)
	if err != nil {
		return nil, err
	}
	dirPathsJson := string(dirPathsJsonBytes)
	if err = PutKV(bucketKey, dirPathsJson); err != nil {
		return nil, err
	}
	return addDirFiles, nil
}

func DeleteBucketDirs(clientModel string, bucket string, paths []string) error {
	bucketKey := FormatBucketDirsKey(bucket)
	var dirPaths []string
	if dirFiles, err := GetBucketDirs(clientModel, bucket); err == nil {
		for _, file := range dirFiles {
			dirPaths = append(dirPaths, file.Path)
		}
	}
	var newDirPaths []string
	for _, dirPath := range dirPaths {
		var delete = false
		for _, path := range paths {
			if dirPath == path {
				delete = true
			}
		}
		if !delete {
			newDirPaths = append(newDirPaths, dirPath)
		}
	}
	dirPathsJsonBytes, err := json.Marshal(newDirPaths)
	if err != nil {
		return err
	}
	dirPathsJson := string(dirPathsJsonBytes)
	if err = PutKV(bucketKey, dirPathsJson); err != nil {
		return err
	}
	return nil
}

func ListBucketFiles(clientModel string, bucket string) ([]File, error) {
	var dirPaths []string
	dirFiles, err := GetBucketDirs(clientModel, bucket)
	if err != nil {
		return nil, err
	}
	for _, file := range dirFiles {
		dirPaths = append(dirPaths, file.Path)
	}
	if len(dirPaths) <= 0 {
		return nil, errors.New("bucket dir invalid")
	}
	var bucketFiles []File
	for _, path := range dirPaths {
		if files, err := ListTypeFiles(path, bucket); err == nil {
			bucketFiles = append(bucketFiles, files...)
		}
	}
	return bucketFiles, nil
}

func FormatBucketDirsKey(bucket string) string {
	return KV_KEY_BUCKET_DIRS + "_" + strings.ToUpper(bucket)
}
