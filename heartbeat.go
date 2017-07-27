package main

import (
	"fmt"
	"github.com/little-cui/sc-tester/helper"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	SERVER_ADDRESS = "127.0.0.1:30100"
	API            = "/registry/v3/microservices/:serviceId/instances/:instanceId/heartbeat"
)

var HEADERS http.Header = http.Header{
	"X-Domain-Name": []string{"default"},
	"Content-Type":  []string{"application/json"},
}

func main() {
	serviceId := helper.GetServiceCenterId()
	instanceId := helper.GetServiceCenterInstanceId(serviceId)
	u := url.URL{
		Scheme: "http",
		Host:   SERVER_ADDRESS,
		Path:   strings.Replace(strings.Replace(API, ":serviceId", serviceId, 1), ":instanceId", instanceId, 1),
	}
	req, err := http.NewRequest("PUT", u.String(), nil)
	req.Header = HEADERS
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.StatusCode, string(body))
}
