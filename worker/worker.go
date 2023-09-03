package worker

import "github.com/lxhcaicai/loj-judge/envexec"

// EnvironmentPool 定义用于执行命令的环境池
type EnvironmentPool interface {
	Get() (envexec.Environment, error)
	Put(envexec.Environment)
}
