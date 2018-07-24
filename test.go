package main

import (
	"fmt"
	"github.com/little-cui/sc-tester/api"
	"time"
)

const (
	CONCURRENT = 50
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

			inst1 := j * 2
			api.RegisterTesterInst(d, j, inst1)
			//go api.Watch(d, j, inst1)

			inst2 := j*2 + 1
			api.RegisterTesterInst(d, j, inst2)
			//go api.Watch(d, j, inst2)

			go loop(d, j, inst2, func(d, s, i int) {
				api.HeartbeatTesterInst(d, s, inst1)
				api.HeartbeatTesterInst(d, s, inst2)
				api.FindTesterInsts(d, s, inst1)
			})
		}
		fmt.Println("Register", max, "OK")
	}()
}

func main() {
	for i := 0; i < CONCURRENT; i++ {
		run(api.SVC_PER_DOMAIN, i)
	}
	<-make(chan struct{})
}
