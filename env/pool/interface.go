package pool

import (
	"github.com/lxhcaicai/loj-judge/envexec"
	"github.com/lxhcaicai/loj-judge/worker"
	"sync"
)

// Environment 定义了 envexec.Environment 的销毁
type Environment interface {
	envexec.Environment
	Reset() error
	Destroy() error
}

// EnvBuilder 定义容器环境的抽象
type EnvBuilder interface {
	Build() (Environment, error)
}

type pool struct {
	builder EnvBuilder

	env []Environment
	mu  sync.Mutex
}

func (p pool) Get() (envexec.Environment, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.env) > 0 {
		rt := p.env[len(p.env)-1]
		p.env = p.env[:len(p.env)-1]
		return rt, nil
	}
	return p.builder.Build()
}

func (p pool) Put(env envexec.Environment) {
	e, ok := env.(Environment)
	if !ok {
		panic("invalid environment put")
	}
	// 如果容器在执行后死亡，不要将其放入池中
	if err := e.Reset(); err != nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.env = append(p.env, e)
}

func NewPool(builder EnvBuilder) worker.EnvironmentPool {
	return &pool{
		builder: builder,
	}
}
