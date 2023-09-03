package envexec

import (
	"context"
	"golang.org/x/sync/errgroup"
)

// Pipe 定义了并行Cmd之间的管道
type Pipe struct {
	// In, Out 定义管道输入源和输出目的地
	In, Out PipeIndex

	// Name 如果名称不为空并且启用了代理，则定义复制出条目名称
	Name string

	// Limit 定义了从代理和代理中复制的最大字节数
	// 超过限制后复制数据
	Limit Size

	// Proxy 创建2个管道，并通过复制数据连接它们
	Proxy bool
}

// PipeIndex 定义命令行的索引和该命令行的fd
type PipeIndex struct {
	Index int
	Fd    int
}

// Group 定义运行指令以运行多个
// 并行执行限制在cGroup内
type Group struct {
	// Cmd定义了在多个环境中并行运行的Cmd
	Cmd []*Cmd

	// Pipes 定义Cmd之间的潜在映射。
	// 确保在相应的CMD中使用nil作为占位符
	Pipes []Pipe

	// NewStoreFile 定义用于创建存储文件的接口
	NewStoreFile NewStoreFile
}

// Run 启动CMD并返回执行结果
func (r *Group) Run(ctx context.Context) ([]Result, error) {
	// prepare files
	fds, pipeToCollect, err := prepareFds(r, r.NewStoreFile)
	if err != nil {
		return nil, err
	}

	// 等待所有CMD命令完成
	var g errgroup.Group
	result := make([]Result, len(r.Cmd))
	for i, c := range r.Cmd {
		i, c := i, c
		g.Go(func() error {
			r, err := runSingle(ctx, c, fds[i], pipeToCollect[i], r.NewStoreFile)
			result[i] = r
			if err != nil {
				result[i].Status = StatusInternalError
				result[i].Error = err.Error()
				return err
			}
			return nil
		})
	}
	err = g.Wait()
	return result, nil
}
