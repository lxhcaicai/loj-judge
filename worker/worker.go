package worker

import (
	"context"
	"github.com/lxhcaicai/loj-judge/envexec"
	"github.com/lxhcaicai/loj-judge/filestore"
	"sync"
	"time"
)

// EnvironmentPool 定义用于执行命令的环境池
type EnvironmentPool interface {
	Get() (envexec.Environment, error)
	Put(envexec.Environment)
}

// Config 定义 worker 配置
type Config struct {
	FileStore             filestore.FileStore
	EnvironmentPool       EnvironmentPool
	Parallelism           int
	WorkDir               string
	TimeLimitTickInterval time.Duration
	ExtraMemoryLimit      envexec.Size
	OutputLimit           envexec.Size
	CopyOutLimit          envexec.Size
	OpenFileLimit         uint64
	ExecObserver          func(Response)
}

// Worker 为执行器定义接口
type Worker interface {
	Start()
	Submit(context.Context, *Request) (<-chan Response, <-chan struct{})
	Execute(context.Context, *Request) <-chan Response
	Shutdown()
}

// worker defines executor worker
type worker struct {
	fs          filestore.FileStore
	envPool     EnvironmentPool
	parallelism int
	workDir     string

	timeLimitTickInterval time.Duration
	extraMemoryLimit      envexec.Size
	outputLimit           envexec.Size
	copyOutLimit          envexec.Size
	openFileLimit         uint64

	execObserver func(Response)

	startOne sync.Once
	stopOne  sync.Once
	wg       sync.WaitGroup
	workCh   chan workRequest
	done     chan struct{}
}

type workRequest struct {
	*Request
	context.Context
	started  chan<- struct{}
	resultCh chan<- Response
}

// New creates new worker
func New(conf Config) Worker {
	return &worker{
		fs:                    conf.FileStore,
		envPool:               conf.EnvironmentPool,
		parallelism:           conf.Parallelism,
		workDir:               conf.WorkDir,
		timeLimitTickInterval: conf.TimeLimitTickInterval,
		extraMemoryLimit:      conf.ExtraMemoryLimit,
		outputLimit:           conf.OutputLimit,
		copyOutLimit:          conf.CopyOutLimit,
		openFileLimit:         conf.OpenFileLimit,
		execObserver:          conf.ExecObserver,
	}
}

func (w worker) Start() {
	//TODO implement me
	panic("implement me")
}

func (w worker) Submit(ctx context.Context, request *Request) (<-chan Response, <-chan struct{}) {
	//TODO implement me
	panic("implement me")
}

func (w worker) Execute(ctx context.Context, request *Request) <-chan Response {
	//TODO implement me
	panic("implement me")
}

func (w worker) Shutdown() {
	//TODO implement me
	panic("implement me")
}
