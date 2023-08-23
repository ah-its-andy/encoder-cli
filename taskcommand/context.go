package taskcommand

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Context struct {
	sources []map[string]string
}

func NewContext() *Context {
	return &Context{sources: make([]map[string]string, 0)}
}

func (c *Context) OutputDir() (string, error) {
	output, ok := c.Get("output")
	if !ok {
		return "", fmt.Errorf("output not set")
	}
	return output, nil
}

func (c *Context) WorkDir() ([]string, error) {
	workdir, ok := c.Get("workdir")
	if !ok {
		return nil, fmt.Errorf("workdir not set")
	}
	cdTempdir, ok := c.Get("cd-tempdir")
	if ok && cdTempdir == "True" {
		var hitFiles []string

		// Walk the directory tree
		err := filepath.WalkDir(workdir, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() && strings.HasSuffix(d.Name(), "_temp") {
				hitFiles = append(hitFiles, d.Name())
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return hitFiles, nil
	} else {
		return []string{workdir}, nil
	}
}

func (c *Context) Clone() *Context {
	nc := NewContext()
	for _, source := range c.sources {
		nc.AddSource(source)
	}
	return nc
}

func (c *Context) AddSource(source map[string]string) {
	c.sources = append(c.sources, source)
}

func (c *Context) UseOptions(opts *CommandOptions) {
	if len(opts.WorkDir) > 0 {
		c.AddSource(map[string]string{
			"workdir": opts.WorkDir,
		})
	}
	if len(opts.OutputDir) > 0 {
		c.AddSource(map[string]string{
			"output": opts.OutputDir,
		})
	}
	if opts.CdTempDir {
		c.AddSource(map[string]string{
			"cd-tempdir": "True",
		})
	}
	if len(opts.Args) > 0 {
		c.AddSource(opts.Args)
	}
}

func (c *Context) getVal(name string) (val string, ok bool) {
	ok = false
	for i := len(c.sources) - 1; i >= 0; i-- {
		if v, o := c.sources[i][name]; o {
			val = v
			ok = true
			break
		}
	}
	return
}

func (c *Context) Get(name string) (val string, ok bool) {
	regular := regexp.MustCompile(`\$\{[a-zA-Z0-9-_]+\}`)
	val, ok = c.get(name)
	if !ok {
		return val, false
	}
	if regular.MatchString(val) {
		return c.get(val)
	}
	return val, true
}

func (c *Context) get(name string) (val string, ok bool) {
	//${a-zA-Z0-9-_}
	regular := regexp.MustCompile(`\$\{[a-zA-Z0-9-_]+\}`)
	matches := regular.FindAllString(name, 0)
	if len(matches) > 0 {
		tmp := map[string]string{}
		retVal := name
		for _, match := range matches {
			if _, ok := tmp[match]; !ok {
				k := match[2 : len(match)-1]
				val, ok := c.getVal(name)
				if !ok {
					log.Printf("variable %s not found, value will be empty", name)
					return retVal, false
				}
				tmp[k] = val
				retVal = strings.ReplaceAll(retVal, match, val)
			}
		}
		return retVal, true
	} else {
		return name, true
	}
}

type CommandOptions struct {
	Name       string            `json:"name" yaml:"name"`
	Command    string            `json:"command" yaml:"command"`
	WorkDir    string            `json:"workdir" yaml:"workdir"`
	OutputDir  string            `json:"outputdir" yaml:"outputdir"`
	SetContext string            `json:"set-context" yaml:"set-context"`
	CdTempDir  bool              `json:"cd-tempdir" yaml:"cd-tempdir"`
	Args       map[string]string `json:"args" yaml:"args"`
}

type TaskList struct {
	Context map[string]string `json:"context" yaml:"context"`
	Tasks   []*CommandOptions `json:"tasks" yaml:"tasks"`
}

func ReadTaskList(fileName string) (*TaskList, error) {
	f, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var ret TaskList
	if filepath.Ext(fileName) == ".json" {
		if err := json.Unmarshal(f, &ret); err != nil {
			return nil, err
		}
	} else if filepath.Ext(fileName) == ".yaml" || filepath.Ext(fileName) == ".yml" {
		if err := yaml.Unmarshal(f, &ret); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("unsupported file extension %s", fileName)
	}
	return &ret, nil
}
