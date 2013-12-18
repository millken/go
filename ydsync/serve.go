package main

import (
	"github.com/qiniu/log"
	"github.com/howeyc/fsnotify"
	"runtime"
	"os"
	"strings"
	"path/filepath"
	"os/exec"
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
var group = make(map[string]string, len(Config.Rsync)) //create a string array with length
var dirmonitor = make(chan *DirMonitor, QUEUE_LENGTH)
var rsync = make(chan *Rsync, QUEUE_LENGTH)

func runMonitor() {
	runtime.GOMAXPROCS(Config.Default.NumCores)
	log.Printf("Started: %d cores, %d threads", Config.Default.NumCores, Config.Default.NumThreads)

	for i := 0; i < Config.Default.NumThreads; i++ {
		go start_worker(dirmonitor)
		go start_rsync(rsync)
	}	
	for k,v := range Config.Rsync {
		group[v.Path] = k
		_, err := fsnotify.NewWatcher()
		if err != nil {
			log.Printf("Failed to watch  %s\n", err)
			return
		}
		log.Infof("config groups : %s, %s", k, v.Path)
		paths = make([]*Paths, 0)
		err = addSubDirs(v.Path)
		if err != nil {
			log.Errorf("subdir error : %s", err)
			return 
		}
	}
	go Inotify()	
	
}

func Inotify() {
	for {
		select {
		case event := <-watchDir.Event:
			//http://code.google.com/p/sersync/source/browse/Inotify.cpp
			fname := filepath.Base(event.Name)
				
			if strings.HasPrefix(fname, ".") || strings.HasSuffix(fname, "~") || fname == "4913" {
				log.Debugf("ignored path: %s", event.Name)
			}else		
			if (event.Mask() & fsnotify.IN_MOVED_FROM) == fsnotify.IN_MOVED_FROM || (event.Mask() & fsnotify.IN_DELETE) == fsnotify.IN_DELETE {
				if event.IsDir() {
					rsync <- &Rsync{"deletedir", event.Name}
				}else{
					rsync <- &Rsync{"deletefile", event.Name}
				}
			}else
			if event.IsCreate() {			
				if event.IsDir() {
					addSubDirs(event.Name)
					rsync <- &Rsync{"createdir", event.Name}
				}else{
					rsync <- &Rsync{"createfile", event.Name}			
				}
			}else
			if  event.IsModify() {		
				if event.IsDir() != true {
					rsync <- &Rsync{"modifyfile", event.Name}
				}

            }

			log.Debugf("serve  event: %s, mask = %d", event, event.Mask())
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
				//log.Debugf("prefix %s: %s", path, v)
			}
		}
		
		//paths = append(paths, &Paths{path})
		//log.Debugf(path)
	}
	return nil
}

func addSubDirs(dirName string)error {
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
		switch r.act{
			case "deletedir": go rsync_deletedir(r.file)
			case "deletefile": go rsync_deletefile(r.file)
			case "createdir" :  go rsync_createdir(r.file)
			case "modifyfile": go rsync_modifyfile(r.file)
			case "createfile": go rsync_modifyfile(r.file)
			default :
		}		
		//log.Debugf("event: %s=>%s", r.act, r.file)
	}
}

//cd /root/test && rsync -artuz --contimeout=3  -R --delete ./   --include="a/" --include="a/f/" --include="a/f/g/" --include="a/f/g/a.txt" --exclude=*  192.168.3.104::test
func rsync_deletedir(file string) {
	log.Debugf("remove watch: %s\n", file)
	watchDir.RemoveWatch(file)
	var name = ""
	for k,v := range group {
		if filepath.HasPrefix(file, k) {
			name = v

		}
	}

	rname := file[len(Config.Rsync[name].Path):]//=substring
	command := "cd " + Config.Rsync[name].Path + " && rsync -artuz --contimeout=3 -R  --delete ./"
	temp := ""
	for _,f := range strings.Split(strings.TrimLeft(rname, "/"), "/") {
		temp = temp + f + "/"
		command += " --include=\"" + temp + "\""

	}
	command += " --include=\"" + temp + "***\"" + " --exclude=* "
	cmd := ""
	for _,host := range Config.Rsync[name].Host {
		cmd = command + host + "::" + name
		log.Debugf("rsync deletedir=> %s", cmd)
		_, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Warnf("command run error : %s", err)
		}			
		
	}
	
}

func rsync_deletefile(file string) {
	var name = ""
	for k,v := range group {
		if filepath.HasPrefix(file, k) {
			name = v

		}
	}

	rname := file[len(Config.Rsync[name].Path):]//=substring
	command := "cd " + Config.Rsync[name].Path + " && rsync -artuz --contimeout=3 "
	temp := ""
	rname1 := strings.TrimLeft(rname, "/")
	if strings.LastIndex(rname1, "/") > -1 {
		rname2 := rname1[:strings.LastIndex(rname1, "/")]
		for _,f := range strings.Split(rname2, "/") {
			temp = temp + f + "/"
			command += " --include=\"" + temp + "\""

		}
	}
	cmd := ""
	for _,host := range Config.Rsync[name].Host {
		cmd = command + " --include=\"" + rname1 + "\" " + host + "::" + name
		log.Debugf("rsync deletefile=> %s", cmd)
		_, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Warnf("command run error : %s", err)
		}		
		
	}
}

func rsync_createdir(file string) {
	var name = ""
	for k,v := range group {
		if filepath.HasPrefix(file, k) {
			name = v

		}
	}

	rname := file[len(Config.Rsync[name].Path):]//=substring
	command := "cd " + Config.Rsync[name].Path + " && rsync -artuz --contimeout=3 -R ./ "

	rname1 := strings.TrimLeft(rname, "/")
	cmd := ""
	for _,host := range Config.Rsync[name].Host {
		cmd = command + " --include=\"./" + rname1 + "\" " + host + "::" + name
		log.Debugf("rsync createdir=> %s", cmd)
		_, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Warnf("command run error : %s", err)
		}		
	}		
}

func rsync_modifyfile(file string) {
	var name = ""
	for k,v := range group {
		if filepath.HasPrefix(file, k) {
			name = v

		}
	}

	rname := file[len(Config.Rsync[name].Path):]//=substring
	command := "cd " + Config.Rsync[name].Path + " && rsync -artuz --contimeout=3 -R ./ "

	rname1 := strings.TrimLeft(rname, "/")
	cmd := ""
	for _,host := range Config.Rsync[name].Host {
		cmd = command + " --include=\"./" + rname1 + "\" " + host + "::" + name
		log.Debugf("rsync modifyfile=> %s", cmd)
		_, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Warnf("command run error : %s", err)
		}		
	}		
}
