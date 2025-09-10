package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"strings"

	"github.com/ulangch/nas_desktop_app/backend/macro"
)

func GetStringMD5(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

func GetFileMD5(path string) string {
	nBytes := 20
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return ""
	}
	if fileInfo.IsDir() {
		return ""
	}
	var bytes []byte
	if fileInfo.Size() <= int64(nBytes)*2 {
		bytes = make([]byte, fileInfo.Size())
		if _, err = file.Read(bytes); err != nil {
			return ""
		}
	} else {
		firstBytes := make([]byte, nBytes)
		if _, err = file.Read(firstBytes); err != nil {
			return ""
		}
		if _, err = file.Seek(-int64(nBytes), io.SeekEnd); err != nil {
			return ""
		}
		lastBytes := make([]byte, nBytes)
		if _, err = file.Read(lastBytes); err != nil {
			return ""
		}
		bytes = append(firstBytes, lastBytes...)
	}
	hash := md5.Sum(bytes)
	md5 := hex.EncodeToString(hash[:])
	return md5
}

/**
 * C:\Users\test\docs       -> /c/Users/test/docs
 * D:\Work\project\src      -> /d/Work/project/src
 */
func ToUnixPath(path string) string {
	// 替换反斜杠为斜杠
	p := strings.ReplaceAll(path, "\\", "/")
	// 处理盘符，比如 "C:/Users/test" → "/c/Users/test"
	if len(p) > 1 && p[1] == ':' {
		drive := strings.ToLower(string(p[0]))
		p = "/" + drive + p[2:] // 去掉 "C:"，前面加 "/c"
	}
	return p
}

func ToPlatformPath(path string) string {
	if macro.IsWin() {
		return ToWindowsPath(path)
	} else {
		return ToUnixPath(path)
	}
}

/**
 * /c/Users/test/docs       -> C:\Users\test\docs
 * /d/Work/project/src      -> D:\Work\project\src
 * /home/user/data          -> \home\user\data
 */
func ToWindowsPath(path string) string {
	p := path
	// 如果是 /c/... 这种格式 → C:\...
	if strings.HasPrefix(p, "/") && len(p) > 2 && p[2] == '/' {
		drive := strings.ToUpper(string(p[1]))
		p = drive + ":" + p[2:]
	}
	// 替换 / 为 \
	p = strings.ReplaceAll(p, "/", `\`)
	return p
}

func UnixRootDir(unixPath string) string {
	path := strings.TrimPrefix(unixPath, "/")
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func IsValidFileName(filename string) bool {
	return len(filename) > 0 && filename[0] != '.'
}
