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
			filiter := 0
			if ev.IsModify() || ev.IsCreate() {
				path, err := os.Lstat(ev.Name)
				if err != nil {
					log.Warnf("os Lstat error: %s", err)
					filiter = 1
				}

				if path.IsDir() != true {
					fname := path.Name()
					if strings.HasPrefix(fname, ".") || strings.HasSuffix(fname, "~") || fname == "4913" {
						filiter = 1
					}
				} 	
			}
			if filiter != 0 {
				log.Warnf("filiter path : %s", ev.Name)
			}else if ev.IsCreate() {
				path, err := os.Lstat(ev.Name)
				if err != nil {
					log.Warnf("os Lstat error: %s", err)
				}				
				if path.IsDir() {
					rsync <- &Rsync{"createdir", ev.Name}
					getSubDirs(ev.Name)
				}else{
					rsync <- &Rsync{"createfile", ev.Name}			
				}
			}else if ev.IsModify() {
				path, err := os.Lstat(ev.Name)
				if err != nil {
					log.Warnf("os Lstat error: %s", err)
				}				
				if path.IsDir() != true {
					fname := path.Name()
					if strings.HasPrefix(fname, ".") || strings.HasSuffix(fname, "~") || fname == "4913" {
						//log.Debugf("filter : %s", ev.Name)
					}else{
						rsync <- &Rsync{"modifyfile", ev.Name}
					}
					//log.Infof("file changed: %s", path.Name())
				}
			}else if ev.IsDelete() {
				if ev.IsDir() {
					rsync <- &Rsync{"deletedir", ev.Name}
				}else{
					rsync <- &Rsync{"deletefile", ev.Name}
				} 
			}else if  ev.IsRename() {
				if path, statErr := os.Stat(ev.Name); os.IsNotExist(statErr) {
					if ev.IsDir() {
						rsync <- &Rsync{"deletedir", ev.Name}
					}else{
						rsync <- &Rsync{"deletefile", ev.Name}
					}
				}else {
				if path.IsDir() {
					rsync <- &Rsync{"createdir", ev.Name}
					getSubDirs(ev.Name)
				}else{
					rsync <- &Rsync{"createfile", ev.Name}			
				}
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
				//log.Debugf("prefix %s: %s", path, v)
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
		//log.Debugf("add watch: %s => %s\n", dirmonitor.Name, dirmonitor.Dir)
		
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
		_, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
		}			
		log.Debugf("rsync deletedir=> %s", cmd)
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
	command := "cd " + Config.Rsync[name].Path + " && rsync -artuz --contimeout=3 -R  --delete ./"
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
		_, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
		}		
		log.Debugf("rsync deletefile=> %s", cmd)
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
		_, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
		}		
		log.Debugf("rsync createdir=> %s", cmd)
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
	command := "cd " + Config.Rsync[name].Path + " && rsync -artuz --contimeout=3 -R  ./ "

	rname1 := strings.TrimLeft(rname, "/")
	cmd := ""
	for _,host := range Config.Rsync[name].Host {
		cmd = command + " --include=\"./" + rname1 + "\" " + host + "::" + name
		_, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
		}		
		log.Debugf("rsync modifyfile=> %s", cmd)
	}		
}