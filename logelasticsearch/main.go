package main

import (
	"flag"
	"github.com/millken/logger"
	"os"
	"os/signal"
)

var (
	VERSION    string = "0.1"
	config     tomlConfig
	configPath string
	gitVersion string
)

func init() {
	if len(gitVersion) > 0 {
		VERSION = VERSION + "/" + gitVersion
	}
}

func main() {
	var err error
	flag.StringVar(&configPath, "config", "ngx2es.toml", "config path")
	flag.Parse()
	logger.Info("Loading config : %s, version: %s", configPath, VERSION)
	err = LoadConfig(configPath)
	if err != nil {
		logger.Exitf("Read config failed.Err = %s", err.Error())
	}
	startKafkaService()
	sigChan := make(chan os.Signal, 3)

	signal.Notify(sigChan, os.Interrupt, os.Kill)

	<-sigChan
}
