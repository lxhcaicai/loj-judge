package linuxcontainer

import (
	"github.com/criyle/go-sandbox/pkg/cgroup"
	"github.com/lxhcaicai/loj-judge/envexec"
	"time"
)

var _ Cgroup = &wCgroup{}

type wCgroup struct {
	cg        cgroup.Cgroup
	cfsPeriod time.Duration
}

func (c *wCgroup) SetCpuset(s string) error {
	return c.cg.SetCPUSet([]byte(s))
}

func (c *wCgroup) SetMemoryLimit(size envexec.Size) error {
	return c.cg.SetMemoryLimit(uint64(size))
}

func (c *wCgroup) SetProcLimit(limit uint64) error {
	return c.cg.SetProcLimit(limit)
}

func (c *wCgroup) SetCPURate(s uint64) error {
	quota := time.Duration(uint64(c.cfsPeriod) * s / 1000)
	return c.cg.SetCPUBandwidth(uint64(quota.Microseconds()), uint64(c.cfsPeriod.Microseconds()))
}

func (c *wCgroup) CPUUsage() (time.Duration, error) {
	t, err := c.cg.CPUUsage()
	return time.Duration(t), err
}

func (c *wCgroup) CurrentMemory() (envexec.Size, error) {
	s, err := c.cg.MemoryUsage()
	return envexec.Size(s), err
}

func (c *wCgroup) MaxMemory() (envexec.Size, error) {
	s, err := c.cg.MemoryUsage()
	return envexec.Size(s), err
}

func (c *wCgroup) AddProc(pid int) error {
	return c.cg.AddProc(pid)
}

func (c *wCgroup) Reset() error {
	return nil
}

func (c *wCgroup) Destory() error {
	return c.cg.Destroy()
}
