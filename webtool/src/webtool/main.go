package main

import (
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

// var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func sigHandler() {
	// TODO On Windows, these signals will not be triggered on closing cmd
	// window. How to detect this?
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGHUP)

	for sig := range sigChan {
		info.Printf("%v caught, exit\n", sig)
		storeSiteStat()
		break
	}
	/*
		if *cpuprofile != "" {
			pprof.StopCPUProfile()
		}
	*/
	os.Exit(0)
}

func main() {
	go sigHandler()
}
