package main

import (
	"context"
	crypto_rand "crypto/rand"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/lxhcaicai/loj-judge/cmd/executorserver/config"
	restexecutor "github.com/lxhcaicai/loj-judge/cmd/executorserver/rest_executor"
	"github.com/lxhcaicai/loj-judge/cmd/executorserver/version"
	"github.com/lxhcaicai/loj-judge/env"
	"github.com/lxhcaicai/loj-judge/env/pool"
	"github.com/lxhcaicai/loj-judge/envexec"
	"github.com/lxhcaicai/loj-judge/filestore"
	"github.com/lxhcaicai/loj-judge/worker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"log"
	math_rand "math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"
)

var logger *zap.Logger

func main() {
	conf := loadConf()
	if conf.Version {
		fmt.Print(version.Version)
		return
	}
	initLogger(conf)
	defer logger.Sync()

	logger.Sugar().Infof("config loaded: %+v", conf)
	initRand()
	warnIfNotLinux()

	// Init environment pool
	fs, _ := newFilesStore(conf)
	b, builderParam := newEnvBuilder(conf)
	envPool := newEnvPool(b, conf.EnableMetrics)
	prefork(envPool, conf.PreFork)
	work := newWorker(conf, envPool, fs)
	work.Start()
	servers := []initFunc{
		initHTTPServer(conf, fs, builderParam),
	}

	// 优雅停机
	sig := make(chan os.Signal, 1+len(servers))

	stops := []stopFunc{}
	for _, s := range servers {
		start, stop := s()
		if start != nil {
			go func() {
				start()
				sig <- os.Interrupt
			}()
		}
		if stop != nil {
			stops = append(stops, stop)
		}
	}

	// 优雅关闭
	signal.Notify(sig, os.Interrupt)
	<-sig
	signal.Reset(os.Interrupt)

	logger.Sugar().Info("Shutting Down...")
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()

	var eg errgroup.Group
	for _, s := range stops {
		s := s
		eg.Go(func() error {
			return s(ctx)
		})
	}

	go func() {
		logger.Sugar().Info("Shutdown Finished", eg.Wait())
		cancel()
	}()

	<-ctx.Done()
}

func loadConf() *config.Config {
	var conf config.Config
	if err := conf.Load(); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		log.Fatalln("load config failed", err)
	}
	return &conf
}

type stopFunc func(ctx context.Context) error
type initFunc func() (start func(), cleanUp stopFunc)

func initHTTPServer(conf *config.Config, fs filestore.FileStore, builderParam map[string]any) initFunc {
	return func() (start func(), cleanUp stopFunc) {
		// Init http handle
		r := initHTTPMux(conf, fs, builderParam)
		srv := http.Server{
			Addr:    conf.HTTPAddr,
			Handler: r,
		}

		return func() {
				lis, err := newListener(conf.HTTPAddr)
				if err != nil {
					logger.Sugar().Error("Http server listen failed: ", err)
					return
				}
				logger.Sugar().Info("Starting http server at ", conf.HTTPAddr, " with listener ", printListener(lis))
				if err := srv.Serve(lis); errors.Is(err, http.ErrServerClosed) {
					logger.Sugar().Info("Http server stopped: ", err)
				} else {
					logger.Sugar().Error("Http server stopped: ", err)
				}
			}, func(ctx context.Context) error {
				logger.Sugar().Info("Http server shutdown")
				return srv.Shutdown(ctx)
			}
	}
}

func initHTTPMux(conf *config.Config, fs filestore.FileStore, builderParam map[string]any) http.Handler {
	var r *gin.Engine
	if conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	r = gin.New()
	r.Use(ginzap.Ginzap(logger, "", false))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	// Version handle
	r.GET("/version", generateHandleVersion(conf))

	// Config handle
	r.GET("/config", generateHandleConfig(conf, builderParam))

	restHandle := restexecutor.New(fs, conf.SrcPrefix, logger)
	restHandle.Register(r)

	return r
}

func generateHandleVersion(conf *config.Config) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"buildVersion":    version.Version,
			"goVersion":       runtime.Version(),
			"platform":        runtime.GOARCH,
			"OS":              runtime.GOOS,
			"copyOutOptional": true,
			"pipeProxy":       true,
			"symlink":         true,
		})
	}
}

func generateHandleConfig(conf *config.Config, builderParam map[string]any) func(context2 *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"copyOutOptional": true,
			"pipeProxy":       true,
			"symlink":         true,
			"fileStorePath":   conf.Dir,
			"runnerConfig":    builderParam,
		})
	}
}

func initLogger(conf *config.Config) {
	if conf.Slient {
		logger = zap.NewNop()
	}

	var err error
	if conf.Release {
		logger, err = zap.NewProduction()
	} else {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		if !conf.EnableDebug {
			config.Level.SetLevel(zap.InfoLevel)
		}
		logger, err = config.Build()
	}
	if err != nil {
		log.Fatalln("init logger failed", err)
	}
}

func initRand() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		logger.Fatal("random generator init failed", zap.Error(err))
	}
	sd := int64(binary.LittleEndian.Uint64(b[:]))
	logger.Sugar().Infof("random seed: %d", sd)
	math_rand.Seed(sd)
}

func warnIfNotLinux() {
	if runtime.GOOS != "linux" {
		logger.Sugar().Warn("Platform is ", runtime.GOOS)
		logger.Sugar().Warn("Please notice that the primary supporting platform is Linux")
		logger.Sugar().Warn("Windows and macOS support are only recommended in development environment")
	}
}

func newFilesStore(conf *config.Config) (filestore.FileStore, func() error) {
	const timeoutCheckInterval = 15 * time.Second
	var cleanUp func() error

	var fs filestore.FileStore
	if conf.Dir == "" {
		if runtime.GOOS == "linux" {
			conf.Dir = "/dev/shm"
		} else {
			conf.Dir = os.TempDir()
		}
		var err error
		conf.Dir, err = os.MkdirTemp(conf.Dir, "executorserver")
		if err != nil {
			logger.Sugar().Fatal("failed to create file store temp dir", err)
		}
		cleanUp = func() error {
			return os.RemoveAll(conf.Dir)
		}
	}
	os.MkdirAll(conf.Dir, 0755)
	fs = filestore.NewFileLocalStore(conf.Dir)
	if conf.EnableDebug {
		fs = newMetricsFileStore(fs)
	}
	if conf.FileTimeout > 0 {
		fs = filestore.NewTimeout(fs, conf.FileTimeout, timeoutCheckInterval)
	}
	return fs, cleanUp
}

func newEnvBuilder(conf *config.Config) (pool.EnvBuilder, map[string]any) {
	b, param, err := env.NewBuilder(env.Config{
		ContainerInitPath:  conf.ContainerInitPath,
		MountConf:          conf.MountConf,
		TmpFsParam:         conf.TmpFsParam,
		NetShare:           conf.NetShare,
		CgroupPrefix:       conf.CgroupPrefix,
		Cpuset:             conf.Cpuset,
		ContainerCredStart: conf.ContainerCredStart,
		EnableCPURate:      conf.EnableCPURate,
		CPUCfsPeriod:       conf.CPUCfsPeriod,
		SeccompConf:        conf.SeccompConf,
		Logger:             logger.Sugar(),
	})
	if err != nil {
		logger.Sugar().Fatal("create environment builder failed", err)
	}
	if conf.EnableMetrics {
		b = &metriceEnvBuilder{b}
	}
	return b, param
}

func newEnvPool(b pool.EnvBuilder, enableMetrics bool) worker.EnvironmentPool {
	p := pool.NewPool(b)
	if enableMetrics {
		p = &metricsEnvPool{p}
	}
	return p
}

func prefork(envPool worker.EnvironmentPool, prefork int) {
	if prefork <= 0 {
		return
	}
	logger.Sugar().Info("create ", prefork, " prefork containers")
	m := make([]envexec.Environment, 0, prefork)
	for i := 0; i < prefork; i++ {
		e, err := envPool.Get()
		if err != nil {
			log.Fatalln("prefork environment failed", err)
		}
		m = append(m, e)
	}
	for _, e := range m {
		envPool.Put(e)
	}
}

func newWorker(conf *config.Config, envPool worker.EnvironmentPool, fs filestore.FileStore) worker.Worker {
	return worker.New(worker.Config{
		FileStore:             fs,
		EnvironmentPool:       envPool,
		Parallelism:           conf.Parallelism,
		WorkDir:               conf.Dir,
		TimeLimitTickInterval: conf.TimeLimitCheckerInterval,
		ExtraMemoryLimit:      *conf.ExtraMemoryLimit,
		OutputLimit:           *conf.OutputLimit,
		CopyOutLimit:          *conf.CopyOutLimit,
		OpenFileLimit:         uint64(conf.OpenFileLimit),
		ExecObserver:          execObserve,
	})
}
