package taskcommand

import (
	"fmt"
	"log"

	"github.com/ah-its-andy/encoder-cli/mkvtoolnix"
	"github.com/ah-its-andy/encoder-cli/utils"
)

func mkvextract(ctx *Context, opts *CommandOptions) error {
	workDirs, err := ctx.WorkDir()
	if err != nil {
		return err
	}
	for _, workDir := range workDirs {
		log.Printf("Walking directory %s", workDir)
		files, err := utils.FilterFiles(workDir)
		if err != nil {
			return err
		}
		for _, f := range files {
			log.Printf("Processing file %s for task %s\r\n", f, opts.Name)
			tempDir, err := utils.TempDir(f)
			if err != nil {
				log.Printf("Error: creating temporary directory for task %s failed, %v\r\n", opts.Name, err)
				return fmt.Errorf("error: processing file failed")
			}
			_, err = mkvtoolnix.ExtractMKV(f, tempDir)
			if err != nil {
				log.Printf("Error: extractMKV failed: %v", err)
				return fmt.Errorf("error: processing file failed")
			}
		}
	}
	return nil
}
