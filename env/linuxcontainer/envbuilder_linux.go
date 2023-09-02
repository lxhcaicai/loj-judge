package linuxcontainer

import (
	"github.com/lxhcaicai/loj-judge/env/pool"
	"syscall"
)

type Config struct {
	Builder    EnvironmentBuilder
	CgroupPool CgroupPool
	WorkDir    string
	Seccomp    []syscall.SockFilter
	Cpuset     string
	CPURate    bool
}

type environmentBuilder struct {
	builder EnvironmentBuilder
	cgPool  CgroupPool
	workDir string
	seccomp []syscall.SockFilter
	cpuset  string
	cpuRate bool
}

func (e environmentBuilder) Build() (pool.Environment, error) {
	//TODO implement me
	panic("implement me")
}

// Build creates linux container
func NewEnvBuilder(c Config) pool.EnvBuilder {
	return &environmentBuilder{
		builder: c.Builder,
		cgPool:  c.CgroupPool,
		workDir: c.WorkDir,
		seccomp: c.Seccomp,
		cpuset:  c.Cpuset,
		cpuRate: c.CPURate,
	}
}
