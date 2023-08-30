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
	"github.com/lxhcaicai/loj-judge/cmd/executorserver/version"
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

	servers := []initFunc{
		initHTTPServer(conf),
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

func initHTTPServer(conf *config.Config) initFunc {
	return func() (start func(), cleanUp stopFunc) {
		// Init http handle
		r := initHTTPMux(conf)
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

func initHTTPMux(conf *config.Config) http.Handler {
	var r *gin.Engine
	if conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	r = gin.New()
	r.Use(ginzap.Ginzap(logger, "", false))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	// Version handle
	r.GET("/version", generateHandleVersion(conf))

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
