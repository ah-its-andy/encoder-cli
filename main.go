package main

import (
	"flag"
	"log"
	"os"

	"github.com/ah-its-andy/encoder-cli/bootstrap"
	"github.com/ah-its-andy/encoder-cli/taskcommand"
)

func main() {
	conf := flag.String("c", "", "")
	task := flag.String("t", "", "")
	flag.Parse()

	if conf == nil || len(*conf) == 0 {
		log.Printf("-c is required")
		os.Exit(-1)
	}

	if task == nil || len(*task) == 0 {
		log.Printf("-t is required")
		os.Exit(-1)
	}

	bootstrap.InitGoConf(*conf)

	taskcommand.RunTask(*task)
}
