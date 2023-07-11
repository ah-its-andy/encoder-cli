package taskcommand

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/ah-its-andy/encoder-cli/mkvtoolnix"
)

func mkvmerge(ctx *Context, opts *CommandOptions) error {
	workDirs, err := ctx.WorkDir()
	if err != nil {
		return err
	}
	outpurDir, err := ctx.OutputDir()
	if err != nil {
		return err
	}
	for _, workDir := range workDirs {
		log.Printf("mkv merging %s\r\n", workDir)
		err = mkvtoolnix.Merge(workDir, filepath.Join(outpurDir, workDir[:len(workDir)-5])+".mkv")
		if err != nil {
			log.Printf("error: merging mkv failed: %v\r\n", err)
			return fmt.Errorf("error: merging mkv failed, dir: %s", workDir)
		}
	}
	return nil
}
