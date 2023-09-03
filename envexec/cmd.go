package envexec

import "github.com/criyle/go-sandbox/runner"

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
