package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/little-cui/sc-tester/helper"
	"net/http"
	"net/url"
	"strings"
)

const (
	LISTWATCH_API = "/registry/v3/microservices/:serviceId/listwatcher"
	WATCH_API     = "/registry/v3/microservices/:serviceId/watcher"
)

var (
	HEADERS http.Header = http.Header{
		"X-Domain-Name": []string{"default"},
	}
	lwCh      chan struct{}
	wCh       chan struct{}
	serviceId string
)

func init() {
	lwCh = make(chan struct{})
	wCh = make(chan struct{})
}

func listwatch() {
	defer close(lwCh)
	fmt.Println("start listwatcher...")

	path := url.URL{
		Scheme: "ws",
		Host:   helper.GetServiceCenterAddress(),
		Path:   strings.Replace(LISTWATCH_API, ":serviceId", serviceId, 1),
	}
	conn, _, err := websocket.DefaultDialer.Dial(path.String(), HEADERS)
	if err != nil {
		panic(err)
	}
	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		if t == websocket.TextMessage {
			fmt.Println("listwatcher:", string(msg))
		}
	}
	conn.Close()
}

func watch() {
	defer close(wCh)
	fmt.Println("start watcher...")

	path := url.URL{
		Scheme: "ws",
		Host:   helper.GetServiceCenterAddress(),
		Path:   strings.Replace(WATCH_API, ":serviceId", serviceId, 1),
	}
	conn, _, err := websocket.DefaultDialer.Dial(path.String(), HEADERS)
	if err != nil {
		panic(err)
	}
	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		if t == websocket.TextMessage {
			fmt.Println("watcher:", string(msg))
		}
	}
	conn.Close()
}

func main() {
	serviceId = helper.GetServiceCenterId()
	go listwatch()
	go watch()

	<-lwCh
	<-wCh
}
