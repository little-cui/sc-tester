package main

import (
	"fmt"
	"github.com/little-cui/sc-tester/helper"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const API = "/registry/v3/microservices/:serviceId/instances"

var HEADERS http.Header = http.Header{
	"X-Domain-Name": []string{"default"},
}

func main() {
	serviceId := helper.GetServiceCenterId()
	HEADERS.Set("X-ConsumerId", serviceId)
	resp, err := http.DefaultClient.Do(&http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   helper.GetServiceCenterAddress(),
			Path:   strings.Replace(API, ":serviceId", serviceId, 1),
		},
		Header: HEADERS,
	})
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}
