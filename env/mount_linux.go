package env

import (
	"fmt"
	"github.com/criyle/go-sandbox/container"
	"github.com/criyle/go-sandbox/pkg/mount"
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

type Mount struct {
	Type     string `yaml:"type"`
	Data     string `yaml:"data"`
	Target   string `yaml:"target"`
	ReadOnly bool   `yaml:"readOnly"`
	Source   string `yaml:"source"`
}

type Mounts struct {
	Mount      []Mount  `yaml:"mount"`
	SymLinks   []Link   `yaml:"symLink"`
	MaskPaths  []string `yaml:"maskPath"`
	WorkDir    string   `yaml:"workDir"`
	HostName   string   `yaml:"hostName"`
	DomainName string   `yaml:"domainName"`
	UID        int      `yaml:"uid"`
	GID        int      `yaml:"gid"`
	Proc       bool     `yaml:"proc"`
}

// Link 定义挂载后要创建的符号链接
type Link struct {
	LinkPath string `yaml:"linkPath"`
	Target   string `yaml:"target"`
}

func readMountConfig(p string) (*Mounts, error) {
	var m Mounts
	d, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(d, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func parseMountConfig(m *Mounts) (*mount.Builder, error) {
	b := mount.NewBuilder()
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for _, mt := range m.Mount {
		target := mt.Target
		if path.IsAbs(target) {
			target = path.Clean(target[1:])
		}
		source := mt.Source
		if !path.IsAbs(source) {
			source = path.Join(wd, source)
		}
		switch mt.Type {
		case "bind":
			b.WithBind(source, target, mt.ReadOnly)
		case "tmpfs":
			b.WithTmpfs(target, mt.Data)
		default:
			return nil, fmt.Errorf("invalid_mount_type: %v", mt.Type)
		}
	}
	if m.Proc {
		b.WithProc()
	}
	return b, nil
}

func getDefaultMount(tmpFsConf string) *mount.Builder {
	return mount.NewBuilder().
		WithBind("/bin", "bin", true).
		WithBind("/lib", "lib", true).
		WithBind("/lib64", "lib64", true).
		WithBind("/usr", "usr", true).
		WithBind("/etc/ld.so.cache", "etc/ld.so.cache", true).
		// Java需要/proc/self/exe作为lib的相对路径
		//然而，/proc给出了类似/proc/1/fd/3的接口。
		//打开该文件将是一个EPERM
		//更改fs的uid和gid将是个好方法
		WithProc().
		// 一些编译器有多个版本
		WithBind("/etc/alternatives", "etc/alternatives", true).
		WithBind("/etc/fpc.cfg", "etc/fpc.cfg", true).
		WithBind("/etc/mono", "etc/mono", true).
		// go wants /dev/null
		WithBind("/dev/null", "dev/null", false).
		// ghc wants /var/lib/ghc
		WithBind("/var/lib/ghc", "var/lib/ghc", true).
		// javaScript wants /dev/urandom
		WithBind("/dev/urandom", "dev/urandom", false).
		// additional devices
		WithBind("/dev/random", "dev/random", false).
		WithBind("/dev/zero", "dev/zero", false).
		WithBind("/dev/full", "dev/full", false).
		// work dir
		WithTmpfs("w", tmpFsConf).
		// tmp dir
		WithTmpfs("tmp", tmpFsConf)
}

var defaultSymLinks = []container.SymbolicLink{
	{LinkPath: "/dev/fd", Target: "/proc/self/fd"},
	{LinkPath: "/dev/stdin", Target: "/proc/self/fd/0"},
	{LinkPath: "/dev/stdout", Target: "/proc/self/fd/1"},
	{LinkPath: "/dev/stderr", Target: "/proc/self/fd/2"},
}

var defaultMaskPaths = []string{
	"/proc/acpi",
	"/proc/asound",
	"/proc/kcore",
	"/proc/keys",
	"/proc/latency_stats",
	"/proc/timer_list",
	"/proc/timer_stats",
	"/proc/sched_debug",
	"/proc/scsi",
	"/usr/lib/wsl/drivers",
	"/usr/lib/wsl/lib",
}
