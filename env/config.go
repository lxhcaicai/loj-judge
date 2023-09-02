package env

import "time"

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type Config struct {
	ContainerInitPath  string
	TmpFsParam         string
	NetShare           bool
	MountConf          string
	SeccompConf        string
	CgroupPrefix       string
	Cpuset             string
	ContainerCredStart int
	EnableCPURate      bool
	CPUCfsPeriod       time.Duration
	Logger
}
