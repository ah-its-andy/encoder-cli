package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"
)

type CommandContext struct {
	cmd *exec.Cmd
	buf *bytes.Buffer
}

func (c *CommandContext) Start() error {
	return c.cmd.Start()
}

func (c *CommandContext) Wait() error {
	return c.cmd.Wait()
}

func (c *CommandContext) Run() error {
	_, err := c.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create StdoutPipe: %v", err)
	}
	if err := c.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}
	if err := c.Wait(); err != nil {
		return fmt.Errorf("command finished with error: %v", err)
	}
	return nil
}

func (c *CommandContext) StdoutPipe() (io.Reader, error) {
	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create StdoutPipe: %v", err)
	}
	go func(reader io.ReadCloser, out io.Writer) {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			c.buf.WriteString(line)
			log.Print(line)
		}
	}(stdout, c.buf)
	return c.buf, nil
}

func Command(name string, arg ...string) *CommandContext {
	cmd := exec.Command(name, arg...)
	return &CommandContext{
		cmd: cmd,
		buf: bytes.NewBuffer(nil),
	}
}

func ExecutableFile(fileName string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s.exe", fileName)
	}
	return fileName
}
