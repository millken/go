package main

import (
	"flag"
	"github.com/millken/logger"
	"runtime"
	"sync"
)

var (
	VERSION    string = "0.1"
	config     tomlConfig
	configPath string
	debugLevel string
	gitVersion string
	workerCh   = make(chan string)
	esCh       = make(chan map[string]interface{})
)

func init() {
	if len(gitVersion) > 0 {
		VERSION = VERSION + "/" + gitVersion
	}
}

func main() {
	var err error
	flag.StringVar(&configPath, "config", "ngx2es.toml", "config path")
	flag.StringVar(&debugLevel, "debug", "INFO", "FINE|DEBUG|TRACE|INFO|ERROR")
	flag.Parse()
	logLevel := logger.INFO
	switch debugLevel {
	case "FINE":
		logLevel = logger.FINE
	case "DEBUG":
		logLevel = logger.DEBUG
	case "TRACE":
		logLevel = logger.TRACE
	case "ERROR":
		logLevel = logger.ERROR
	}
	logger.Global = logger.NewDefaultLogger(logLevel)
	logger.Info("Loading config : %s, version: %s", configPath, VERSION)
	err = LoadConfig(configPath)
	if err != nil {
		logger.Exitf("Read config failed.Err = %s", err.Error())
	}

	numCpus := runtime.NumCPU()
	runtime.GOMAXPROCS(numCpus)

	if err = loadIpdb(); err != nil {
		logger.Exitf("Load ipdb failed.Err = %s", err.Error())
	}
	startElasticSearchService()

	wg := new(sync.WaitGroup)
	for i := 0; i < numCpus; i++ {
		wg.Add(1)
		go startWorker(i, workerCh, wg)
	}

	startKafkaService()

	close(workerCh)
	wg.Wait()

}
