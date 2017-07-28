package main

import (
	"fmt"
	"github.com/little-cui/sc-tester/helper"
	"io/ioutil"
	"net/http"
	"net/url"
)

const API = "/health"

var HEADERS http.Header = http.Header{
	"X-Domain-Name": []string{"default"},
}

func main() {
	resp, err := http.DefaultClient.Do(&http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   helper.GetServiceCenterAddress(),
			Path:   API,
		},
		// Header: HEADERS,
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
