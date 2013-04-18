package main

import (
	"github.com/qiniu/log"
	"github.com/howeyc/fsnotify"
	"runtime"
	"os"
	"path/filepath"
)
const (
	// How many requests/responses can be in the queue before blocking.
	QUEUE_LENGTH = 65535
)
type DirMonitor struct {
	Name string //rsync name
	Dir string
}

type Paths struct {
	path string
}
var paths = make([]*Paths, 0)
var watchDir, err = fsnotify.NewWatcher()


func runMonitor() {
	runtime.GOMAXPROCS(Config.Default.NumCores)
	log.Printf("Started: %d cores, %d threads", Config.Default.NumCores, Config.Default.NumThreads)

	dirmonitor := make(chan *DirMonitor, QUEUE_LENGTH)

	for i := 0; i < Config.Default.NumThreads; i++ {
		go start_worker(dirmonitor)
	}	
	for k,v := range Config.Groups {

		_, err := fsnotify.NewWatcher()
		if err != nil {
			log.Printf("Failed to watch  %s\n", err)
			return
		}
		log.Infof("config groups : %s, %s", k, v.Path)
		paths = make([]*Paths, 0)
		err = getSubDirs(v.Path)
		if err != nil {
			log.Errorf("subdir error : %s", err)
			return 
		}
		for _,p := range paths {
			dirmonitor <- &DirMonitor{k, p.path}
		}
	}
	for {
		select {
		case ev := <-watchDir.Event:
			log.Debugf("serve  event: %s", ev)
		case err := <-watchDir.Error:
			log.Println("fsnotify error:", err)
		}
	}	
}

func WalkFunc(path string, info os.FileInfo, err error)error {
	if info.IsDir() {
		paths = append(paths, &Paths{path})
		//log.Debugf(path)
	}
	return nil
}

func getSubDirs(dirName string)error {
	_, err := os.Stat(dirName)
	if err != nil {
		return err
	}
	err = filepath.Walk(dirName, WalkFunc)
	if err != nil {return err}
	return nil
}

// and returns the results to the results channel.
func start_worker(queue <-chan *DirMonitor) {
	for dirmonitor := range queue {
		watchDir.Watch(dirmonitor.Dir)
		log.Debugf("dirmonitor: %s => %s\n", dirmonitor.Name, dirmonitor.Dir)
		
	}
}

