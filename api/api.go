package api

import (
	"net/url"
	"strings"
	"net/http"
	"time"
	"fmt"
	"github.com/little-cui/sc-tester/helper"
	"io/ioutil"
	"math/rand"
)

const (
	HEARTBEAT_API = "/registry/v3/microservices/:serviceId/instances/:instanceId/heartbeat"
	CREATE_API = "/registry/v3/microservices"
	REGISTER_API = "/registry/v3/microservices/:serviceId/instances"
	INSTANCE_API = "/registry/v3/instances"
	EXIST_API = "/registry/v3/existence"
)

func Create() {
	appId, serviceName, version := "Tester", "TestService", fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	r := strings.NewReader(fmt.Sprintf(`{
	"service":{
		"serviceId":"",
		"appId":"%s",
		"serviceName":"%s",
		"version":"%s",
		"level":"BACK",
		"status":"UP",
		"properties":{
			"allowCrossApp":"true"
		}
	}
}`, appId, serviceName, version))

	client := http.Client{}
	req, err := http.NewRequest("GET", (&url.URL{
		Scheme: "http",
		Host: helper.GetServiceCenterAddress(),
		Path: EXIST_API,
		RawQuery: "type=microservice&appId="+appId+"&serviceName="+serviceName+"&version="+version,
	}).String(), nil)
	req.Header = http.Header{
		"X-Domain-Name": []string{"default"},
		"Content-Type":  []string{"application/json"},
	}
	t := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == http.StatusOK {
		if  time.Now().Sub(t) > time.Second {
			fmt.Println("exist:", string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
		}
		return
	}

	req, err = http.NewRequest("POST", (&url.URL{
		Scheme: "http",
		Host:   helper.GetServiceCenterAddress(),
		Path:   CREATE_API,
	}).String(), r)
	req.Header = http.Header{
		"X-Domain-Name": []string{"default"},
		"Content-Type":  []string{"application/json"},
	}
	t = time.Now()
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK || time.Now().Sub(t) > time.Second {
		fmt.Println("create:", string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
	}
}

func Register() {
	serviceId := helper.GetServiceCenterId()
	r := strings.NewReader(fmt.Sprintf(`{
	"instance": {
		"serviceId":"%s",
		"endpoints":["rest://127.0.0.%d:30100"],
		"hostName":"service_center_10_229_33_15",
		"status":"UP",
		"healthCheck":{"mode":"push","interval":2,"times":3},
		"stage":"prod"
	}
}`, serviceId, rand.Intn(255)))
	u := url.URL{
		Scheme: "http",
		Host:   helper.GetServiceCenterAddress(),
		Path:   strings.Replace(REGISTER_API, ":serviceId", serviceId, 1),
	}
	req, err := http.NewRequest("POST", u.String(), r)
	req.Header = http.Header{
		"X-Domain-Name": []string{"default"},
		"Content-Type":  []string{"application/json"},
	}

	t := time.Now()
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK || time.Now().Sub(t) > time.Second {
		fmt.Println("register:", string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
	}
}

func Heartbeat() {
	serviceId := helper.GetServiceCenterId()
	instanceId := helper.GetServiceCenterInstanceId(serviceId)
	u := url.URL{
		Scheme: "http",
		Host:   helper.GetServiceCenterAddress(),
		Path:   strings.Replace(strings.Replace(HEARTBEAT_API, ":serviceId", serviceId, 1), ":instanceId", instanceId, 1),
	}
	req, err := http.NewRequest("PUT", u.String(), nil)
	req.Header = http.Header{
		"X-Domain-Name": []string{"default"},
		"Content-Type":  []string{"application/json"},
	}

	t := time.Now()
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK || time.Now().Sub(t) > time.Second {
		fmt.Println("heartbeat:", string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
	}
	// fmt.Println(resp.StatusCode, string(body), time.Now().Sub(t))
}

func Find() {
	serviceId := helper.GetServiceCenterId()
	t := time.Now()
	client := http.Client{}
	resp, err := client.Do(&http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   helper.GetServiceCenterAddress(),
			Path:   INSTANCE_API,
			RawQuery: "appId=Tester&serviceName=TestService&version=latest",
		},
		Header: http.Header{
			"X-Domain-Name": []string{"default"},
			"X-ConsumerId": []string{serviceId},
		},
	})
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK || time.Now().Sub(t) > time.Second {
		fmt.Println("find:", string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
	}
}