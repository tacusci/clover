package utils

import (
	"runtime"
	"strings"
)

//TranslatePath to translate file location paths cross OS
func TranslatePath(path string) string {
	if runtime.GOOS == "windows" {
		path = strings.Replace(path, ":/", ":\\", -1)
		path = strings.Replace(path, "/", "\\", -1)
	} else if runtime.GOOS == "darwin" {
		path = strings.Replace(path, "\\", "/", -1)
	}
	return path
}
