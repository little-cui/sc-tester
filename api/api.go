package api

import (
	"fmt"
	"github.com/little-cui/sc-tester/helper"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	HEARTBEAT_API = "/registry/v3/microservices/:serviceId/instances/:instanceId/heartbeat"
	CREATE_API    = "/registry/v3/microservices"
	REGISTER_API  = "/registry/v3/microservices/:serviceId/instances"
	FIND_API      = "/registry/v3/instances"
	EXIST_API     = "/registry/v3/existence"
	INSTANCE_API  = "/registry/v3/microservices/:serviceId/instances/:instanceId"
)

func print(code, timeout bool, args ...interface{}) {
	if code {
		fmt.Fprintln(os.Stderr, args...)
		return
	}
	if timeout {
		fmt.Println(args...)
	}
}

func CreateTesterService() {
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
		Scheme:   "http",
		Host:     helper.GetServiceCenterAddress(),
		Path:     EXIST_API,
		RawQuery: "type=microservice&appId=" + appId + "&serviceName=" + serviceName + "&version=" + version,
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

	print(false, time.Now().Sub(t) > time.Second,
		"exist:", string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
	if resp.StatusCode == http.StatusOK {
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
	print(resp.StatusCode != http.StatusOK, time.Now().Sub(t) > time.Second,
		"Create:", string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
}

func RegisterSCInst() {
	serviceId := helper.GetServiceCenterId()
	r := strings.NewReader(fmt.Sprintf(`{
	"instance": {
		"serviceId":"%s",
		"endpoints":["rest://127.0.0.%d:30100"],
		"hostName":"service_center_10_229_33_15",
		"status":"UP",
		"healthCheck":{"mode":"push","interval":1,"times":3},
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
	print(resp.StatusCode != http.StatusOK, time.Now().Sub(t) > time.Second,
		"Register:", string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
}

func HeartbeatSCInst() {
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
	print(resp.StatusCode != http.StatusOK, time.Now().Sub(t) > time.Second,
		"HeartbeatSCInst:", instanceId, string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
}

func FindTesterInsts() {
	serviceId := helper.GetServiceCenterId()

	versionRules := []string{
		"latest",
		"0+",
		fmt.Sprintf("%d.%d.%d.%d",
			rand.Intn(128), rand.Intn(128), rand.Intn(128), rand.Intn(128)) +
			"-" +
			fmt.Sprintf("%d.%d.%d.%d",
				rand.Intn(128), rand.Intn(128), rand.Intn(128), rand.Intn(128)),
		fmt.Sprintf("%d.%d.%d.%d",
			rand.Intn(128), rand.Intn(128), rand.Intn(128), rand.Intn(128)),
	}
	v := versionRules[rand.Intn(len(versionRules))]

	t := time.Now()
	client := http.Client{}
	resp, err := client.Do(&http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme:   "http",
			Host:     helper.GetServiceCenterAddress(),
			Path:     FIND_API,
			RawQuery: "appId=Tester&serviceName=TestService&version=" + url.QueryEscape(v),
		},
		Header: http.Header{
			"X-Domain-Name": []string{"default"},
			"X-ConsumerId":  []string{serviceId},
		},
	})
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	print(resp.StatusCode != http.StatusOK, time.Now().Sub(t) > time.Second,
		"FindTesterInsts:", v, string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
}

func GetSCInsts() {
	serviceId := helper.GetServiceCenterId()
	t := time.Now()
	client := http.Client{}
	resp, err := client.Do(&http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   helper.GetServiceCenterAddress(),
			Path:   strings.Replace(REGISTER_API, ":serviceId", serviceId, 1),
		},
		Header: http.Header{
			"X-Domain-Name": []string{"default"},
			"X-ConsumerId":  []string{serviceId},
		},
	})
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	print(resp.StatusCode != http.StatusOK, time.Now().Sub(t) > time.Second,
		"GetSCInsts:", string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
}
