package main

import (
	"github.com/little-cui/sc-tester/api"
	"time"
)

const (
	INTERVAL   = 100 * time.Millisecond
	CONCURRENT = 100
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
		run(api.RegisterTesterInst)
		run(api.HeartbeatTesterInst)
		run(api.FindTesterInsts)
		run(api.GetSCInsts)
	}
	<-make(chan struct{})
}
