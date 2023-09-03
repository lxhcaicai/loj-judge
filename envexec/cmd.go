package envexec

import (
	"context"
	"github.com/criyle/go-sandbox/runner"
	"os"
	"time"
)

// Size 表示以字节为单位的数据大小
type Size = runner.Size

// RunnerResult 表示进程完成结果
type RunnerResult = runner.Result

type CmdCopyOutFile struct {
	Name     string // Name 输出到copyOut的文件
	Optional bool   // Optional 如果文件不存在，则忽略该文件
}

type FileErrorType int

type FileError struct {
	Name    string        `json:"name"`
	Type    FileErrorType `json:"type"`
	Message string        `json:"message,omitempty"`
}

// Cmd 定义在容器环境中运行程序的指令
type Cmd struct {
	Environment Environment

	// 在执行前要复制的文件内容
	CopyIn map[string]File

	// 在执行之前创建的符号链接
	SymLinks map[string]string

	// exec argument, environment
	Args []string
	Env  []string

	// Files for the executing command
	Files []File
	TTY   bool // use pty as input / output

	// 资源限制
	TimeLimit         time.Duration
	MemoryLimit       Size
	StackLimit        Size
	ExtraMemoryLimit  Size
	OutputLimit       Size
	ProcLimit         uint64
	OpenFileLimit     uint64
	CPURateLimit      uint64
	StrictMemoryLimit bool
	CpuSetLimit       string

	// 在cmd启动后调用Waiter，它应该返回
	// 一旦超过时间限制。
	// 作为TLE返回true，作为正常退出返回false(上下文结束)
	Waiter func(context.Context, Process) bool

	// 在执行后要复制的文件名
	CopyOut    []CmdCopyOutFile
	CopyOutMax Size

	// 指定转储所有/w内容的目录
	CopyOutDir string
}

type Result struct {
	Status Status

	ExitStatus int

	Error string // error

	Time    time.Duration
	RunTime time.Duration
	Memory  Size // byte

	// Files 存储复制文件
	Files map[string]*os.File

	// 存储文件错误详细信息
	FileError []FileError
}

const (
	ErrCopyInOpenFile FileErrorType = iota
	ErrCopyInCreateDir
	ErrCopyInCreateFile
	ErrCopyInCopyContent
	ErrCopyOutOpen
	ErrCopyOutNotRegularFile
	ErrCopyOutSizeExceeded
	ErrCopyOutCreateFile
	ErrCopyOutCopyContent
	ErrCollectSizeExceeded
	ErrSymlink
)
