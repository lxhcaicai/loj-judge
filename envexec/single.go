package envexec

import "context"

// Single 定义要单独运行的运行指令
// 执行被限制在cgroup
type Single struct {
	// 定义Cmd在多个环境中并行运行
	Cmd *Cmd

	// NewStoreFile 定义用于创建存储文件的接口
	NewStoreFile NewStoreFile
}

// 启动命令行并返回执行结果
func (s *Single) Run(ctx context.Context) (result Result, err error) {
	// prepare files
	fd, pipeToCollect, err := prepareCmdFd(s.Cmd, len(s.Cmd.Files), s.NewStoreFile)
	if err != nil {
		return result, err
	}

	result, err = runSingle(ctx, s.Cmd, fd, pipeToCollect, s.NewStoreFile)
	if err != nil {
		result.Status = StatusInternalError
		result.Error = err.Error()
		return result, err
	}

	return result, err
}
