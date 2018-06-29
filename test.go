package main

import (
	"fmt"
	"github.com/little-cui/sc-tester/api"
	"time"
)

const (
	CONCURRENT = 100
)

func loop(d, s, i int, f func(d, s, i int)) {
	for {
		<-time.After(30 * time.Second)
		f(d, s, i)
	}
}

func run(max, d int) {
	go func() {
		for j := 0; j < max; j++ {
			api.CreateTesterService(d, j)
			inst := j * 2
			api.RegisterTesterInst(d, j, inst)
			go loop(d, j, inst, func(d, s, i int) {
				api.HeartbeatTesterInst(d, s, i)
				api.FindTesterInsts(d, s, i)
			})
			inst = j*2 + 1
			api.RegisterTesterInst(d, j, inst)
			go loop(d, j, inst, api.HeartbeatTesterInst)
		}
		fmt.Println("Register", max, "OK")
	}()
}

func main() {
	fmt.Println("")
	for i := 0; i < CONCURRENT; i++ {
		run(api.SVC_PER_DOMAIN, i)
	}
	<-make(chan struct{})
}
