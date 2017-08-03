package main

import (
	"github.com/little-cui/sc-tester/api"
	"time"
)

const (
	INTERVAL   = 100 * time.Millisecond
	CONCURRENT = 5
)

func run(f func()) {
	go func() {
		for {
			f()
			<-time.After(INTERVAL)
		}
	}()
}

func main() {
	for i := 0; i < CONCURRENT; i++ {
		run(api.CreateTesterService)
		run(api.RegisterSCInst)
		run(api.HeartbeatSCInst)
		run(api.FindTesterInsts)
		run(api.GetSCInsts)
	}
	<-make(chan struct{})
}
