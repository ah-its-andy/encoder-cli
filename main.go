package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/ah-its-andy/encoder-cli/bootstrap"
	"github.com/ah-its-andy/encoder-cli/taskcommand"
)

func main() {
	taskInfoConfig := flag.String("c", "", "")
	flag.Parse()

	if taskInfoConfig == nil || len(*taskInfoConfig) == 0 {
		log.Printf("-c is required")
		os.Exit(-1)
	}

	execPath, _ := filepath.Abs("./")
	bootstrap.InitGoConf(execPath)

	taskcommand.RunTask(*taskInfoConfig)
}
