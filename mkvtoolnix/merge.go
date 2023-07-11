package mkvtoolnix

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ah-its-andy/encoder-cli/utils"
)

func Merge(tempDir string, mkvFilePath string) error {
	fileList, err := ioutil.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	args := []string{"-o", mkvFilePath}
	for _, file := range fileList {
		if !file.IsDir() {
			filePath := filepath.Join(tempDir, file.Name())
			args = append(args, filePath)
		}
	}

	cmd := utils.Command("mkvmerge", args...)
	_, err = cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create StdoutPipe: %v", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command finished with error: %v", err)
	}
	return nil
}
