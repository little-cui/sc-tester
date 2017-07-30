package main

import (
	"fmt"
	"github.com/little-cui/sc-tester/helper"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const API = "/registry/v3/microservices"

var HEADERS http.Header = http.Header{
	"X-Domain-Name": []string{"default"},
	"Content-Type":  []string{"application/json"},
}

func main() {
	r := strings.NewReader(`{
	"service":{
		"serviceId":"",
		"appId":"Tester",
		"serviceName":"TestService",
		"version":"1.0.0",
		"level":"BACK",
		"status":"UP",
		"properties":{
			"allowCrossApp":"true"
		}
	}
}`)
	u := url.URL{
		Scheme: "http",
		Host:   helper.GetServiceCenterAddress(),
		Path:   API,
	}
	req, err := http.NewRequest("POST", u.String(), r)
	req.Header = HEADERS
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}
