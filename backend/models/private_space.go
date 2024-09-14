package models

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const TOKEN_ALIVE_DURATION = 7 * 24 * time.Hour
const TOKEN_UPDATE_DURATION = 1 * 24 * time.Hour

type PrivateSpace struct {
	Sid   string `json:"sid"`
	Token string `json:"token"`
	File  File   `json:"file"`
}

func CreatePrivateSpace(password string) (PrivateSpace, error) {
	sid, _ := GetKV(KV_KEY_PRIVATE_SPACE_SID)
	if sid != "" {
		return PrivateSpace{}, fmt.Errorf("space exists, sid=%s", sid)
	}
	dirFile, err := GetPrivateSpaceDir()
	if err != nil {
		return PrivateSpace{}, err
	}
	sid = uuid.NewString()
	token := EncryptToken(sid, password)
	PutKV(KV_KEY_PRIVATE_SPACE_SID, sid)
	PutKV(KV_KEY_PRIVATE_SPACE_PASSWORD, password)
	PutKV(KV_KEY_PRIVATE_SPACE_TOKEN, token)
	return PrivateSpace{Sid: sid, Token: token, File: dirFile}, nil
}

func GetPrivateSpace(sid string, password string) (PrivateSpace, error) {
	if existSid, _ := GetKV(KV_KEY_PRIVATE_SPACE_SID); sid != existSid {
		return PrivateSpace{}, fmt.Errorf("space not exist, sid=%s", sid)
	}
	if existPassword, _ := GetKV(KV_KEY_PRIVATE_SPACE_PASSWORD); existPassword != password {
		return PrivateSpace{}, errors.New("password invalid")
	}
	encryptToken, _ := GetKV(KV_KEY_PRIVATE_SPACE_TOKEN)
	if encryptToken == "" {
		encryptToken = EncryptToken(sid, password)
		PutKV(KV_KEY_PRIVATE_SPACE_TOKEN, encryptToken)
	}
	dirFile, err := GetPrivateSpaceDir()
	if err != nil {
		return PrivateSpace{}, err
	} else {
		return PrivateSpace{Sid: sid, Token: encryptToken, File: dirFile}, nil
	}
}

func GetPrivateSpaceSid() string {
	sid, _ := GetKV(KV_KEY_PRIVATE_SPACE_SID)
	return sid
}

func GetPrivateSpaceDir() (File, error) {
	disks, err := GetDeviceDisks()
	if err != nil {
		return File{}, errors.New("no device disks")
	}
	dirPath := filepath.Join(disks[0].Path, ".private")
	dirFile, err := CreateDirectoryPerm(dirPath, os.ModePerm)
	if err != nil {
		return File{}, errors.New("make space dir failed")
	} else {
		return dirFile, nil
	}
}

func ValidateToken(encryptToken string) error {
	if existEncryptToken, err := GetKV(KV_KEY_PRIVATE_SPACE_TOKEN); err != nil || existEncryptToken != encryptToken {
		return errors.New("token invalid")
	} else {
		return nil
	}
}

func EncryptToken(sid string, password string) string {
	decryptToken := sid + "##" + password
	return base64.StdEncoding.EncodeToString([]byte(decryptToken))
}

func DecryptToken(encryptToken string) (string, string, error) {
	decryptToken, err := base64.StdEncoding.DecodeString(encryptToken)
	if err != nil {
		return "", "", errors.New("decrypt token failed")
	}
	segments := strings.Split(string(decryptToken), "##")
	if len(segments) != 2 {
		return "", "", errors.New("decrypt token invalid")
	}
	return segments[0], segments[1], nil
}
