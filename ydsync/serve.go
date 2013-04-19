package main

import (
	"github.com/qiniu/log"
	"github.com/howeyc/fsnotify"
	"runtime"
	"os"
	"strings"
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

type Rsync struct {
	act string
	file string
}
var paths = make([]*Paths, 0)
var watchDir, err = fsnotify.NewWatcher()
var group = make(map[string]string, len(Config.Groups)) //create a string array with length
var dirmonitor = make(chan *DirMonitor, QUEUE_LENGTH)
var rsync = make(chan *Rsync, QUEUE_LENGTH)

func runMonitor() {
	runtime.GOMAXPROCS(Config.Default.NumCores)
	log.Printf("Started: %d cores, %d threads", Config.Default.NumCores, Config.Default.NumThreads)

	for i := 0; i < Config.Default.NumThreads; i++ {
		go start_worker(dirmonitor)
		go start_rsync(rsync)
	}	
	for k,v := range Config.Groups {
		group[v.Path] = k
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
		/*
		for _,p := range paths {
			dirmonitor <- &DirMonitor{k, p.path}
		}
		*/
	}
	go Inotify()	
	
}

func Inotify() {
	for {
		select {
		case ev := <-watchDir.Event:
			if ev.IsCreate() {
				path, err := os.Lstat(ev.Name)
				if err != nil {
					log.Warnf("watch error: %s", err)
				}
				if path.IsDir() {
					getSubDirs(ev.Name)
					rsync <- &Rsync{"createdir", ev.Name}
				}
			}else if ev.IsModify() {
				path, err := os.Lstat(ev.Name)
				if err != nil {
					log.Warnf("os Lstat error: %s", err)
				}				
				if path.IsDir() != true {
					fname := path.Name()
					if strings.HasPrefix(fname, ".") || strings.HasSuffix(fname, "~") || fname == "4913" {
						log.Debugf("filter : %s", ev.Name)
					}else{
						rsync <- &Rsync{"modifyfile", ev.Name}
					}
					log.Infof("file changed: %s", path.Name())
				}
			}else if ev.IsDelete() {
				if ev.IsDir() {
					rsync <- &Rsync{"deletedir", ev.Name}
				}else{
					rsync <- &Rsync{"deletefile", ev.Name}
				} 
			}
			log.Debugf("serve  event: %s", ev)
		case err := <-watchDir.Error:
			log.Println("fsnotify error:", err)
		}
	}
}
func WalkFunc(path string, info os.FileInfo, err error)error {
	if info.IsDir() {
		for k,v := range group {
			if filepath.HasPrefix(path, k) {
				dirmonitor <- &DirMonitor{v, path}				
				log.Debugf("prefix %s: %s", path, v)
			}
		}
		
		//paths = append(paths, &Paths{path})
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
		log.Debugf("add watch: %s => %s\n", dirmonitor.Name, dirmonitor.Dir)
		
	}
}

func start_rsync(queue <-chan *Rsync) {
	for r := range queue {
		log.Debugf("rsync file: %s=>%s", r.act, r.file)
	}
}
