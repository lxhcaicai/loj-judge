package pool

// Environment 定义了 envexec.Environment 的销毁
type Environment interface {
}

// EnvBuilder 定义容器环境的抽象
type EnvBuilder interface {
	Build() (Environment, error)
}
