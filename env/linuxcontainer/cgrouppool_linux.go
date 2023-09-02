package linuxcontainer

import (
	"github.com/lxhcaicai/loj-judge/envexec"
	"time"
)

type Cgroup interface {
	SetCpuset(string) error
	SetMemoryLimit(envexec.Size) error
	SetProcLimit(uint642 uint64) error
	SetCPURate(uint642 uint64) error // 1000 as 1

	CPUUsage() (time.Duration, error)
	CurrentMemory() (envexec.Size, error)
	MaxMemory() (envexec.Size, error)

	AddProc(int) error
	Reset() error
	Destory() error
}

// CgroupPool implements pool of Cgroup
type CgroupPool interface {
	Get() (Cgroup, error)
	Put(Cgroup)
}
