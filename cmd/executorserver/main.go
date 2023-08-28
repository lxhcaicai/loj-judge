package main

import (
	"flag"
	"fmt"
	"github.com/lxhcaicai/loj-judge/cmd/executorserver/config"
	"github.com/lxhcaicai/loj-judge/cmd/executorserver/version"
	"log"
	"os"
)

func main() {
	conf := loadConf()
	if !conf.Version {
		fmt.Print(version.Version)
		return
	}
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
