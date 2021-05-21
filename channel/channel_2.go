package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var reqs []int
	for i := 1; i < 150; i++ {
		reqs = append(reqs, i)
	}
	wg.Add(len(reqs))
	semaphore := make(chan struct{}, 5)
	for i := 0; i < len(reqs); i++ {
		go func(i int) {
			semaphore <- struct{}{}
			run(reqs[i])

			wg.Done()
			func() { <-semaphore }()
		}(i)
	}
	wg.Wait()
}

func run(i int) {
	log.Println(i)
	time.Sleep(time.Second * time.Duration(rand.Intn(10)))
	return
}
