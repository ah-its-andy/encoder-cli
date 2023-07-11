package utils

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func TempDir(fileName string) (string, error) {
	tempDir, err := ioutil.TempDir(filepath.Dir(fileName), filepath.Base(fileName[:len(fileName)-len(filepath.Ext(fileName))]+"_temp"))
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %v", err)
	}
	return tempDir, nil
}
