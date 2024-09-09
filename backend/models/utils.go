package models

import "path/filepath"

var imageExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg", ".ico", ".heic"}

var videoExtensions = []string{".mp4", ".avi", ".mkv", ".mov", ".flv", ".wmv", ".webm", ".mpg", ".mpeg", ".m4v", ".3gp", ".mts"}

func isImage(name string) bool {
	extention := filepath.Ext(name)
	for _, ext := range imageExtensions {
		if ext == extention {
			return true
		}
	}
	return false
}

func isVideo(name string) bool {
	extention := filepath.Ext(name)
	for _, ext := range videoExtensions {
		if ext == extention {
			return true
		}
	}
	return false
}
