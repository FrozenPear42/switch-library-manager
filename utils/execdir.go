package utils

import (
	"os"
	"path/filepath"
)

func GetExecDir() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	dirPath := filepath.Dir(path)
	return dirPath, nil
}
