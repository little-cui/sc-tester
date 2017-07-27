package main

import (
	"fmt"
	"github.com/little-cui/sc-tester/helper"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const API = "/registry/v3/microservices/:serviceId/instances/:instanceId/heartbeat"

var HEADERS http.Header = http.Header{
	"X-Domain-Name": []string{"default"},
	"Content-Type":  []string{"application/json"},
}

func main() {
	serviceId := helper.GetServiceCenterId()
	instanceId := helper.GetServiceCenterInstanceId(serviceId)
	u := url.URL{
		Scheme: "http",
		Host:   helper.GetServiceCenterAddress(),
		Path:   strings.Replace(strings.Replace(API, ":serviceId", serviceId, 1), ":instanceId", instanceId, 1),
	}
	req, err := http.NewRequest("PUT", u.String(), nil)
	req.Header = HEADERS

	t := time.Now()
	resp, err := http.DefaultClient.Do(req)
	fmt.Println("spend:", time.Now().Sub(t))
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.StatusCode, string(body), time.Now().Sub(t))
}
