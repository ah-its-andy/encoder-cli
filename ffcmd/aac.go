package ffcmd

import (
	"fmt"
	"path/filepath"

	"github.com/ah-its-andy/encoder-cli/utils"
	"github.com/ah-its-andy/goconf"
)

func ConvertToAAC(fileName string, outputFilePath string) error {
	// Get the file extension
	ext := filepath.Ext(fileName)

	// Check if the file is already in AAC format
	if ext == ".aac" {
		return nil
	}

	execPath, ok := goconf.GetString("tools.ffmpeg")
	if !ok {
		return fmt.Errorf("ffmpeg not found")
	}
	// Run the ffmpeg command to convert the audio file to AAC format
	cmd := utils.Command(filepath.Join(execPath, utils.ExecutableFile("ffmpeg")), "-i", fileName, "-c:a", "aac", "-b:a", "192k", "-strict", "-2", outputFilePath)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
