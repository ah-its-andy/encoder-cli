package taskcommand

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ah-its-andy/encoder-cli/ffcmd"
	"github.com/ah-its-andy/encoder-cli/utils"
	"github.com/ah-its-andy/goconf"
)

func toaac(ctx *Context, opts *CommandOptions) error {
	workDirs, err := ctx.WorkDir()
	if err != nil {
		return err
	}
	exts := []string{".truehd", ".dts", ".dtshd", ".mp3", ".aac3"}
	extsRaw, ok := goconf.GetSection("toaac.exts").GetRaw()
	if ok {
		if val, ok := extsRaw.([]any); ok {
			exts = make([]string, 0)
			for _, v := range val {
				exts = append(exts, v.(string))
			}
		}
	}
	for _, workDir := range workDirs {
		files, err := utils.FilterFiles(workDir, exts...)
		if err != nil {
			return err
		}
		for _, file := range files {
			dir := filepath.Dir(file)
			ext := filepath.Ext(file)
			aacFileName := fmt.Sprintf("%s.aac", file[len(dir):len(file)-len(ext)])
			if _, err := os.Stat(aacFileName); os.IsExist(err) {
				log.Printf("aac audio file already exists, skipping %s\r\n", aacFileName)
				continue
			}
			if err := ffcmd.ConvertToAAC(file, filepath.Join(workDir, aacFileName)); err != nil {
				log.Printf("error: converting file to aac failed: %v\r\n", err)
				return fmt.Errorf("error: converting file to aac failed")
			}
		}
	}
	return nil
}
