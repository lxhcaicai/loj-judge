package config

import (
	"github.com/koding/multiconfig"
	"github.com/lxhcaicai/loj-judge/envexec"
	"os"
	"time"
)

type Config struct {

	// container
	ContainerInitPath  string `flagUsage:"container init path"`
	PreFork            int    `flagUsage:"control # of the prefork workers" default:"1"`
	TmpFsParam         string `flagUsage:"tmpfs mount data (only for default mount with no mount.yaml)" default:"size=128m,nr_inodes=4k"`
	NetShare           bool   `flagUsage:"share net namespace with host"`
	MountConf          string `flagUsage:"specifies mount configuration file" default:"mount.yaml"`
	SeccompConf        string `flagUsage:"specifies seccomp filter" default:"seccomp.yaml"`
	Parallelism        int    `flagUsage:"control the # of concurrency execution (default equal to number of cpu)"`
	CgroupPrefix       string `flagUsage:"control cgroup prefix" default:"executor_server"`
	ContainerCredStart int    `flagUsage:"control the start uid&gid for container (0 uses unprivileged root)" default:"0"`

	// file store
	SrcPrefix []string `flagUsage:"specifies directory prefix for source type copyin (example: -src-prefix=/home,/usr)"`
	Dir       string   `flagUsage:"specifies directory to store file upload / download (in memory by default)"`

	// runner limit
	FileTimeout              time.Duration `flagUsage:"specified timeout for filestore files"`
	Cpuset                   string        `flagUsage:"control the usage of cpuset for all containerd process"`
	CPUCfsPeriod             time.Duration `flagUsage:"set cpu.cfs_period" default:"100ms"`
	EnableCPURate            bool          `flagUsage:"enable cpu cgroup rate control"`
	TimeLimitCheckerInterval time.Duration `flagUsage:"specifies time limit checker interval" default:"100ms"`
	ExtraMemoryLimit         *envexec.Size `flagUsage:"specifies extra memory buffer for check memory limit" default:"16k"`
	OutputLimit              *envexec.Size `flagUsage:"specifies POSIX rlimit for output for each command" default:"256m"`
	CopyOutLimit             *envexec.Size `flagUsage:"specifies default file copy out max" default:"256m"`
	OpenFileLimit            int           `flagUsage:"specifies max open file count" default:"256"`

	// server config
	HTTPAddr      string `flagUsage:"specifies the http binding address"`
	EnableDebug   bool   `flagUsage:"enable debug endpoint"`
	EnableMetrics bool   `flagUsage:"enable promethus metrics endpoint"`

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
