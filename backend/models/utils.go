package models

import (
	"path/filepath"
	"strings"
)

const FILE_TYPE_IMAGE = "image"
const FILE_TYPE_VIDEO = "video"
const FILE_TYPE_AUDIO = "audio"
const FILE_TYPE_DOC = "doc"
const FILE_TYPE_NOVEL = "novel"
const FILE_TYPE_ZIP = "zip"

var imageExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg", ".ico", ".heic"}
var videoExtensions = []string{".mp4", ".avi", ".mkv", ".mov", ".flv", ".wmv", ".webm", ".mpg", ".mpeg", ".m4v", ".3gp", ".mts"}
var audioExtensions = []string{".mp3", ".wav", ".flac", ".aac", ".ogg", ".wma", ".m4a", ".aiff", ".alac"}
var documentExtensions = []string{".txt", ".md", ".rtf", ".odt", ".pdf", ".ppt", ".ai", ".psd", ".doc", ".docx", ".xls", ".xml"}
var novelExtensions = []string{".txt"}
var zipExtension = []string{".zip", ".rar"}

func IsImage(name string) bool {
	extention := filepath.Ext(name)
	for _, ext := range imageExtensions {
		if ext == extention {
			return true
		}
	}
	return false
}

func IsVideo(name string) bool {
	extention := filepath.Ext(name)
	for _, ext := range videoExtensions {
		if ext == extention {
			return true
		}
	}
	return false
}

func IsAudio(name string) bool {
	extention := filepath.Ext(name)
	for _, ext := range audioExtensions {
		if ext == extention {
			return true
		}
	}
	return false
}

func IsDocument(name string) bool {
	extention := filepath.Ext(name)
	for _, ext := range documentExtensions {
		if ext == extention {
			return true
		}
	}
	return false
}

func IsNovel(name string) bool {
	extention := filepath.Ext(name)
	for _, ext := range novelExtensions {
		if ext == extention {
			return true
		}
	}
	return false
}

func IsZipOrRar(name string) bool {
	extention := filepath.Ext(name)
	for _, ext := range zipExtension {
		if ext == extention {
			return true
		}
	}
	return false
}

func FilterFileByType(name string, filterType string) bool {
	lowerFilterType := strings.ToLower(filterType)
	if lowerFilterType == FILE_TYPE_IMAGE {
		return IsImage(name)
	} else if lowerFilterType == FILE_TYPE_VIDEO {
		return IsVideo(name)
	} else if lowerFilterType == FILE_TYPE_AUDIO {
		return IsAudio(name)
	} else if lowerFilterType == FILE_TYPE_DOC {
		return IsDocument(name)
	} else if lowerFilterType == FILE_TYPE_NOVEL {
		return IsNovel(name)
	} else if lowerFilterType == FILE_TYPE_ZIP {
		return IsZipOrRar(name)
	}
	return false
}

func FilterFile(name string) bool {
	if len(name) <= 0 || name[0] == '.' {
		return false
	}
	return true
}

const KV_KEY_ID = "KV_KEY_ID"
const KV_KEY_NAME = "KV_KEY_NAME"
const KV_KEY_DISKS = "KV_KEY_DISKS"
const KV_KEY_GALLERY_DIR = "KV_KEY_GALLERY_DIR"
const KV_KEY_UPLOAD_DIR = "KV_KEY_UPLOAD_DIR"
const KV_KEY_BUCKET_DIRS = "KV_KEY_BUCKET_DIRS"
