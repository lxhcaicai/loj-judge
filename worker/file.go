package worker

import (
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
