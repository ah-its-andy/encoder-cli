package bootstrap

import (
	"path/filepath"

	"github.com/ah-its-andy/goconf"
	physicalfile "github.com/ah-its-andy/goconf/physicalFile"
)

func InitGoConf(execPath string) {
	goconf.Init(func(b goconf.Builder) {
		b.AddSource(physicalfile.Yaml(filepath.Join(execPath, "conf", "config.yaml")))
		b.AddSource(physicalfile.Yaml(filepath.Join(execPath, "conf", "tools.yaml")))
	})
}
