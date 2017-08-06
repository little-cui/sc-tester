package main

import (
	"fmt"
	"github.com/little-cui/sc-tester/api"
	"time"
)

const (
	INTERVAL   = 100 * time.Millisecond
	CONCURRENT = 10
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
	t := time.Now()
	for i := 0; i < 6000; i++ {
		api.CreateTesterService()
	}
	fmt.Println("Register 6000 micro-service spend", time.Now().Sub(t))
	t = time.Now()
	for i := 0; i < 50; i++ {
		api.RegisterSCInst()
	}
	fmt.Println("Register 50 micro-service instances spend", time.Now().Sub(t))
	for i := 0; i < CONCURRENT; i++ {
		//run(api.CreateTesterService)
		//run(api.RegisterSCInst)
		run(api.HeartbeatSCInst)
		//run(api.FindTesterInsts)
		//run(api.GetSCInsts)
	}
	<-make(chan struct{})
}
