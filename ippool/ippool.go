package ippool

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

var ErrNoResult = errors.New("empty result")

// Target
type Target struct {
	Host string
	Port int

	Counter  int
	Interval time.Duration
	Timeout  time.Duration
}

func (target Target) String() string {
	return fmt.Sprintf("%s:%d", target.Host, target.Port)
}

type Result struct {
	Counter        int
	SuccessCounter int
	Target         *Target

	MinDuration   time.Duration
	MaxDuration   time.Duration
	TotalDuration time.Duration
}

// Avg return the average time of ping
func (result *Result) Avg() time.Duration {
	if result.SuccessCounter == 0 {
		return 0
	}
	return result.TotalDuration / time.Duration(result.SuccessCounter)
}

func (result *Result) String() string {
	return `
--- ` + result.Target.String() + ` ping statistics ---
` + strconv.Itoa(result.Counter) + ` probes sent, ` + strconv.Itoa(result.SuccessCounter) + ` successful, ` + strconv.Itoa(result.Failed()) + ` failed.
rtt min/avg/max = ` + result.MinDuration.String() + `/` + result.Avg().String() + `/` + result.MaxDuration.String()
}

// Failed return failed counter
func (result Result) Failed() int {
	return result.Counter - result.SuccessCounter
}

type ipPool struct {
	sync.RWMutex
	stop     chan bool
	targets  []*Target
	results  []*Result
	interval time.Duration
}

func NewIPPool(interval time.Duration) *ipPool {
	return &ipPool{
		stop:     make(chan bool, 1),
		interval: interval,
	}
}

func (ip *ipPool) AddTarget(target *Target) {
	ip.Lock()
	defer ip.Unlock()
	ip.targets = append(ip.targets, target)
}

func (ip *ipPool) Reset() {
	ip.Lock()
	defer ip.Unlock()
	ip.targets = nil
}

func (ip *ipPool) Start(ctx context.Context) error {
	go func() {
		t := time.NewTicker(1)
		defer t.Stop()
		for {
			select {
			case <-ip.stop:
				return
			case <-ctx.Done():
				return
			case <-t.C:
				ip.pinger()
				t.Reset(ip.interval)
			}
		}
	}()

	return nil
}

func (ip *ipPool) Results() []*Result {
	ip.RLock()
	defer ip.RUnlock()
	results := ip.results
	return results
}

func (ip *ipPool) Targets() []*Target {
	ip.RLock()
	defer ip.RUnlock()
	targets := ip.targets
	return targets
}

func (ip *ipPool) pinger() {
	resultChan := make(chan *Result)
	for _, target := range ip.Targets() {
		go ip.ping(target, resultChan)
	}
	var results []*Result
	for range ip.targets {
		result := <-resultChan
		results = append(results, result)
		log.Println(result)
	}
	ip.Lock()
	ip.results = results
	ip.Unlock()
}

func (ip *ipPool) ping(target *Target, ch chan<- *Result) {
	result := &Result{
		Target: target,
	}
	t := time.NewTicker(target.Interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if result.Counter >= target.Counter && target.Counter != 0 {
				ch <- result
				return
			}
			dur, err := timeIt(func() error {
				conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", target.Host, target.Port), target.Timeout)
				if err != nil {
					return err
				}
				conn.Close()
				return nil
			})
			result.Counter++
			duration := time.Duration(dur)
			if err == nil {
				if result.MinDuration == 0 {
					result.MinDuration = duration
				}
				if result.MaxDuration == 0 {
					result.MaxDuration = duration
				}
				result.SuccessCounter++
				if duration > result.MaxDuration {
					result.MaxDuration = duration
				} else if duration < result.MinDuration {
					result.MinDuration = duration
				}
				result.TotalDuration += duration
			}
		}
	}

}

func (ip *ipPool) Stop(ctx context.Context) error {
	ip.stop <- true
	return nil
}

func (ip *ipPool) GetFastestTarget() (*Target, error) {
	results := ip.Results()
	if len(results) == 0 {
		return nil, ErrNoResult
	}
	res := results[0]
	for _, result := range results {
		if result.SuccessCounter > result.Failed() && res.Avg() <= result.Avg() {
			res = result
		}
	}
	return res.Target, nil
}

func (ip *ipPool) GetRandomTarget() (*Target, error) {
	results := make([]*Result, 0)

	for _, result := range ip.Results() {
		if result.SuccessCounter > result.Failed() {
			results = append(results, result)
		}
	}
	if len(results) == 0 {
		return nil, ErrNoResult
	}
	rand.Seed(time.Now().UnixNano())
	return results[rand.Intn(len(results))].Target, nil
}
