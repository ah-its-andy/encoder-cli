package bootstrap

import (
	"github.com/ah-its-andy/goconf"
)

func InitGoConf(execPath string) {
	goconf.Init(func(b goconf.Builder) {
		b.AddSource(goconf.Memory(map[string]string{
			"tools.mkvtoolnix": "/usr/bin",
			"tools.ffmpeg":     "/usr/bin",
		}))
		b.AddSource(goconf.EnvironmentVariable(""))
	})
}
