package main

import (
	"code.google.com/p/gcfg"
	"fmt"
	"github.com/howeyc/fsnotify"
	"github.com/qiniu/log"
	"os"
	"path/filepath"
	"time"
)

type AppConfig struct {
	Default struct {
        Debug bool
        NumCores int
        NumThreads int
    }
	
    Groups map[string]*struct {
        Path string
    }
}

var Config = new(AppConfig)

func configWatcher(fileName string) {

	configReader(fileName)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Failed to watch  %s\n", err)
		return
	}

	if err := watcher.Watch(*flagconfig); err != nil {
		fmt.Println(err)
		return
	}

	for {
		select {
		case ev := <-watcher.Event:
			//log.Debugf("config file change ,event: %s", ev)
			if filepath.Clean(ev.Name) == fileName {
				if ev.IsModify() { //only modify event
					time.Sleep(200 * time.Millisecond)
					configReader(fileName)
					log.Debugf("reload config file: %s", fileName)
				}
			}
		case err := <-watcher.Error:
			log.Println("fsnotify error:", err)
		}
	}

}

var lastReadConfig time.Time

func configReader(fileName string) error {

	stat, err := os.Stat(fileName)
	if err != nil {
		log.Printf("Failed to find config file: %s\n", err)
		return err
	}

	if !stat.ModTime().After(lastReadConfig) {
		return err
	}

	lastReadConfig = time.Now()

	log.Printf("Loading config: %s\n", fileName)

	cfg := new(AppConfig)

	err = gcfg.ReadFileInto(cfg, fileName)
	if err != nil {
		log.Printf("Failed to parse config data: %s\n", err)
		return err
	}

	log.Printf("%d\n", len(cfg.Groups))
	log.Println("DEBUG :", cfg.Default.Debug)
	// log.Println("STATHAT FLAG  :", cfg.Flags.HasStatHat)

	Config = cfg

	go runMonitor()
	return nil
}

