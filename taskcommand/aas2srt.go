package taskcommand

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ah-its-andy/encoder-cli/subtitle"
	"github.com/ah-its-andy/encoder-cli/utils"
)

func aastosrt(ctx *Context, opts *CommandOptions) error {
	workDirs, err := ctx.WorkDir()
	if err != nil {
		return err
	}

	for _, workDir := range workDirs {
		files, err := utils.FilterFiles(workDir, ".aas")
		if err != nil {
			return err
		}
		for _, file := range files {
			dir := filepath.Dir(file)
			ext := filepath.Ext(file)
			srtFileName := fmt.Sprintf("%s.srt", file[len(dir):len(file)-len(ext)])
			if _, err := os.Stat(srtFileName); os.IsExist(err) {
				log.Printf("srt file already exists, skipping %s\r\n", srtFileName)
				continue
			}
			if err := subtitle.ConvertAAS2SRT(file, filepath.Join(workDir, srtFileName)); err != nil {
				log.Printf("error: converting file to aac failed: %v\r\n", err)
				return fmt.Errorf("error: converting file to aac failed")
			}
		}
	}
	return nil
}
