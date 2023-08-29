package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lxhcaicai/loj-judge/cmd/executorserver/config"
	"github.com/lxhcaicai/loj-judge/cmd/executorserver/version"
	//"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"
)

//var logger *zap.Logger

func main() {
	conf := loadConf()
	if conf.Version {
		fmt.Print(version.Version)
		return
	}

	//defer logger.Sync()
	//logger.Sugar().Infof("config loaded: %+v", conf)
	servers := []initFunc{
		initHTTPServer(conf),
	}

	// 优雅停机
	sig := make(chan os.Signal, 1+len(servers))

	//
	signal.Notify(sig, os.Interrupt)
	<-sig

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

	//logger.Sugar().Info("Shutting Down...")
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()

	var eg errgroup.Group
	for _, s := range stops {
		s := s
		eg.Go(func() error {
			return s(ctx)
		})
	}

	go func() {
		//logger.Sugar().Info("Shutdown Finished", eg.Wait())
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
					//logger.Sugar().Error("Http server listen failed: ", err)
					return
				}
				// //logger.Sugar().Info("Starting http server at ", conf.HTTPAddr, " with listener ", printListener(lis))
				if err := srv.Serve(lis); errors.Is(err, http.ErrServerClosed) {
					//logger.Sugar().Info("Http server stopped: ", err)
				} else {
					//logger.Sugar().Error("Http server stopped: ", err)
				}
			}, func(ctx context.Context) error {
				//logger.Sugar().Info("Http server shutdown")
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
	//r.Use(ginzap.Ginzap(logger, "", false))
	//r.Use(ginzap.RecoveryWithZap(logger, true))

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
