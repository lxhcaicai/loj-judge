package config

import (
	"github.com/koding/multiconfig"
	"os"
	"time"
)

type Config struct {

	// file store
	Dir string `flagUsage:"flagUsage:"specifies directory to store file upload / download (in memory by default)"`

	// runner limit
	FileTimeout time.Duration `flagUsage:"specified timeout for filestore files"`

	// server config
	HTTPAddr    string `flagUsage:"specifies the http binding address"`
	EnableDebug bool   `flagUsage:"enable debug endpoint"`

	// logger config
	Release bool `flagUsage:"release level of logs"`
	Slient  bool `flagUsage:"do not print logs"`

	// 展示版本并推出
	Version bool `flagUsage:"show version and exit"`
}

// 从标志和环境变量加载配置
func (c *Config) Load() error {
	cl := multiconfig.MultiLoader(
		&multiconfig.TagLoader{},
		&multiconfig.EnvironmentLoader{
			Prefix:    "ES",
			CamelCase: true,
		},
		&multiconfig.FlagLoader{
			CamelCase: true,
			EnvPrefix: "ES",
		},
	)
	if os.Getpid() == 1 {
		c.Release = true
		c.HTTPAddr = ":6060"
	} else {
		c.HTTPAddr = "localhost:6060"
	}
	return cl.Load(c)
}
