package main

import (
	"fmt"
	"github.com/little-cui/sc-tester/api"
	"time"
)

const (
	CONCURRENT = 100
)

func loop(d int, s int, i int) {
	for {
		<-time.After(30 * time.Second)
		api.HeartbeatTesterInst(d, s, i)
	}
}

func run(max, d int) {
	go func() {
		for j := 0; j < max; j++ {
			api.CreateTesterService(d, j)
			inst := j * 2
			api.RegisterTesterInst(d, j, inst)
			go loop(d, j, inst)
			inst = j*2 + 1
			api.RegisterTesterInst(d, j, inst)
			go loop(d, j, inst)
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
