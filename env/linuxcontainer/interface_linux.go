package linuxcontainer

import (
	"github.com/criyle/go-sandbox/container"
	"github.com/criyle/go-sandbox/pkg/cgroup"
)

type EnvironmentBuilder interface {
	Build() (container.Environment, error)
}

// CgroupBuilder 为runner构建cgroup
type CgroupBuilder interface {
	Random(string) (cg cgroup.Cgroup, err error)
}
