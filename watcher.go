package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/little-cui/sc-tester/helper"
	"net/http"
	"net/url"
	"strings"
	"time"
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
	fmt.Printf("start listwatch %s...\n", serviceId)

	path := url.URL{
		Scheme: "ws",
		Host:   helper.GetServiceCenterAddress(),
		Path:   strings.Replace(LISTWATCH_API, ":serviceId", serviceId, 1),
	}
	conn, resp, err := websocket.DefaultDialer.Dial(path.String(), HEADERS)
	if err != nil {
		fmt.Println(resp, err)
		panic(err)
	}
	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("listwatcher:", "error:", err)
			break
		}
		if t == websocket.TextMessage {
			fmt.Println("listwatcher:", string(msg))
		}
		conn.WriteControl(websocket.PingMessage, []byte("sss"), time.Now().Add(10*time.Second))
		go func() {
			<-time.After(4 * time.Second)
			conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(10*time.Second))
		}()
	}
	conn.Close()
}

func watch() {
	defer close(wCh)
	fmt.Printf("start watch %s...\n", serviceId)

	path := url.URL{
		Scheme: "ws",
		Host:   helper.GetServiceCenterAddress(),
		Path:   strings.Replace(WATCH_API, ":serviceId", serviceId, 1),
	}
	conn, resp, err := websocket.DefaultDialer.Dial(path.String(), HEADERS)
	if err != nil {
		fmt.Println(resp, err)
		panic(err)
	}
	i := 0
	for {
		i++
		t, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("watcher:", "error:", err)
			break
		}
		if t == websocket.TextMessage {
			fmt.Println("watcher:", i, string(msg))
		}
		go func() {
			<-time.After(3 * time.Second)
			conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(10*time.Second))
		}()
	}
	conn.Close()
}

func main() {
	serviceId = "0c22c09a804711e8873cfa163e17c38b"
	go listwatch()
	go watch()
	go listwatch()
	go watch()
	go listwatch()
	go watch()

	<-lwCh
	<-wCh
}
