package taskcommand

import (
	"fmt"
	"log"
	"os"
)

var commands = map[string]func(ctx *Context, opts *CommandOptions) error{
	"mkvmerge":   mkvmerge,
	"mkvextract": mkvextract,
	"toaac":      toaac,
	"aastosrt":   aastosrt,
}

func RunTask(configFile string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("error: panic %v\r\n", err)
		}
	}()
	if err := runTask(configFile); err != nil {
		log.Printf("error: %v\r\n", err)
		os.Exit(-1)
	}
}

func runTask(configFile string) error {
	taskList, err := ReadTaskList(configFile)
	if err != nil {
		return fmt.Errorf("reading task list failed: %v", err)
	}
	ctx := NewContext()
	ctx.AddSource(taskList.Context)
	for k, task := range taskList.Tasks {
		task.Name = k
		cmd, ok := commands[task.Command]
		if !ok {
			return fmt.Errorf("unknown command: %v", task.Command)
		}
		scopeCtx := ctx.Clone()
		scopeCtx.UseOptions(task)
		if err := cmd(scopeCtx, task); err != nil {
			return err
		}
	}
	return nil
}
