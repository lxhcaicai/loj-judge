package worker

import (
	"github.com/lxhcaicai/loj-judge/envexec"
	"os"
	"time"
)

type Size = envexec.Size
type CmdCopyOutFile = envexec.CmdCopyOutFile
type PipeMap = envexec.Pipe

// Cmd 定义了在envexec中使用的启动程序的命令和限制
type Cmd struct {
	Args  []string
	Env   []string
	Files []CmdFile
	TTY   bool

	CPULimit          time.Duration
	ClockLimit        time.Duration
	MemoryLimit       Size
	StackLimit        Size
	OutputLimit       Size
	ProcLimit         uint64
	OpenFileLimit     uint64
	CPURateLimit      uint64
	CPUSetLimit       string
	StrictMemoryLimit bool

	CopyIn   map[string]CmdFile
	Symlinks map[string]string

	CopyOut       []CmdCopyOutFile
	CopyOutCached []CmdCopyOutFile
	CopyOutMax    uint64
	CopyOutDir    string
}

// Request 定义单个worker请求
type Request struct {
	RequestID   string
	Cmd         []Cmd
	PipeMapping []PipeMap
}

// Result 定义单个命令响应
type Result struct {
	Status     envexec.Status
	ExitStatus int
	Error      string
	Time       time.Duration
	RunTime    time.Duration
	Memory     envexec.Size
	Files      map[string]*os.File
	FileIDs    map[string]string
	FileError  []envexec.FileError
}

// Response 定义单个请求的工作响应
type Response struct {
	RequestID string
	Results   []Result
	Error     error
}
