package utils

import (
	"os"
	"path/filepath"
)

func FilterFiles(dirPath string, exts ...string) ([]string, error) {
	var hitFiles []string

	// Walk the directory tree
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		// Check if the file extension is .mkv
		for _, ext := range exts {
			if filepath.Ext(path) == ext {
				hitFiles = append(hitFiles, path)
				break
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return hitFiles, nil
}
