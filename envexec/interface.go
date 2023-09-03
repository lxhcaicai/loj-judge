package envexec

import (
	"context"
	"os"
	"time"
)

// ExecveParam 是在环境中运行进程的参数
type ExecveParam struct {
	Args []string

	Env []string

	Files []uintptr

	ExecFile uintptr

	TTY bool

	Limit Limit
}

// Limit 定义进程运行的资源限制
type Limit struct {
	Time         time.Duration // Time limit
	Memory       Size          // Memory limit
	Proc         uint64        // Process count limit
	Stack        Size          // Stack limit
	Output       Size          // Output limit
	Rate         uint64        // CPU Rate limit
	OpenFile     uint64        // Number of open files
	CPUSet       string        // CPU set limit
	StrictMemory bool          // Use stricter memory limit (e.g. rlimit)
}

// Usage 定义峰值进程资源使用情况
type Usage struct {
	Time   time.Duration
	Memory Size
}

// Process 正在运行的进程组的进程引用
type Process interface {
	Done() <-chan struct{} // Done 返回一个等待进程退出的通道
	Result() RunnerResult  // Result 等待完成并返回RunnerResult
	Usage() Usage          // Usage检索运行时的资源占用情况
}

// Environment defines the interface to access container execution environment
type Environment interface {
	Execve(context.Context, ExecveParam) (Process, error)
	WorkDir() *os.File // WorkDir 返回打开的工作目录，之后不应关闭
	// Open 使用给定的相对路径和标志打开工作目录下的文件
	Open(path string, flags int, perm os.FileMode) (*os.File, error)
	// 在容器内创建目录
	MkdirAll(path string, perm os.FileMode) error
	// 为文件/目录创建符号链接
	Symlink(oldName, newName string) error
}

type NewStoreFile func() (*os.File, error)
