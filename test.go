package main

import (
	"time"
	"github.com/little-cui/sc-tester/api"
)

const (
	INTERVAL = 200*time.Millisecond
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
	for i:= 0;i<CONCURRENT;i++ {
	run(api.Create)
	run(api.Register)
	run(api.Heartbeat)
	run(api.Find)
	}
	<-make(chan struct {})
}
