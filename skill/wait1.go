//http://stackoverflow.com/questions/18405023/how-would-you-define-a-pool-of-goroutines-to-be-executed-at-once-in-golang
package main

import (
	"os/exec"
	"strconv"
	"sync"
)

func main() {
	tasks := make(chan *exec.Cmd, 64)

	// spawn four worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			for cmd := range tasks {
				cmd.Run()
			}
			wg.Done()
		}()
	}

	// generate some tasks
	for i := 0; i < 10; i++ {
		tasks <- exec.Command("zenity", "--info", "--text='Hello from iteration n."+strconv.Itoa(i)+"'")
	}
	close(tasks)

	// wait for the workers to finish
	wg.Wait()
}
