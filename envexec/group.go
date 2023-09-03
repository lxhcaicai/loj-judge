package envexec

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
