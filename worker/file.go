package worker

import (
	"bytes"
	"fmt"
	"github.com/lxhcaicai/loj-judge/envexec"
	"github.com/lxhcaicai/loj-judge/filestore"
)

// CmdFile 定义命令行中使用的文件
type CmdFile interface {
	// EnvFile 为环境准备文件
	EnvFile(fs filestore.FileStore) (envexec.File, error)
	// 用于打印调试信息的字符串
	String() string
}

var (
	_ CmdFile = &LocalFile{}
	_ CmdFile = &MemoryFile{}
	_ CmdFile = &CachedFile{}
	_ CmdFile = &Collector{}
)

// LocalFile 定义本地文件系统上的文件存储
type LocalFile struct {
	Src string
}

// LocalFile 为envexec文件准备文件
func (f *LocalFile) EnvFile(fs filestore.FileStore) (envexec.File, error) {
	return envexec.NewFileInput(f.Src), nil
}

func (f *LocalFile) String() string {
	return fmt.Sprintf("local: %s", f.Src)
}

// MemoryFile 定义内存中的文件存储
type MemoryFile struct {
	Content []byte
}

// EnvFile 为envexec文件准备文件
func (f *MemoryFile) EnvFile(fs filestore.FileStore) (envexec.File, error) {
	return envexec.NewFileReader(bytes.NewReader(f.Content), false), nil
}

func (f *MemoryFile) String() string {
	return fmt.Sprintf("memory:(len:%d)", len(f.Content))
}

// CachedFile 定义缓存在文件存储中的文件
type CachedFile struct {
	FileID string
}

// EnvFile 为envexec文件准备文件
func (f *CachedFile) EnvFile(fs filestore.FileStore) (envexec.File, error) {
	_, fd := fs.Get(f.FileID)
	if fd == nil {
		return nil, fmt.Errorf("file not exists with id %v", f.FileID)
	}
	return fd, nil
}

func (f *CachedFile) String() string {
	return fmt.Sprintf("cached:(fileId:%s)", f.FileID)
}

// Collector 在要通过管道收集的输出(stdout / stderr)上定义
type Collector struct {
	Name string       // 生成copyOut的伪名称
	Max  envexec.Size // 需要收集的大小
	Pipe bool
}

// EnvFile 为envexec文件准备文件
func (f *Collector) EnvFile(fs filestore.FileStore) (envexec.File, error) {
	return envexec.NewFileCollector(f.Name, f.Max, f.Pipe), nil
}

func (f *Collector) String() string {
	return fmt.Sprintf("collector:(name:%s, max:%d, pipe:%v)", f.Name, f.Max, f.Pipe)
}
