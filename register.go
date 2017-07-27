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
	"Content-Type":  []string{"application/json"},
}

func main() {
	serviceId := helper.GetServiceCenterId()
	r := strings.NewReader(fmt.Sprintf(`{
	"instance": {
		"serviceId":"%s",
		"endpoints":["rest://127.0.0.2:30100"],
		"hostName":"service_center_10_229_33_15",
		"status":"UP",
		"healthCheck":{"mode":"push","interval":2,"times":3},
		"stage":"prod"
	}
}`, serviceId))
	u := url.URL{
		Scheme: "http",
		Host:   helper.GetServiceCenterAddress(),
		Path:   strings.Replace(API, ":serviceId", serviceId, 1),
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
