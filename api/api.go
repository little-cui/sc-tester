package api

import (
	"fmt"
	"github.com/little-cui/sc-tester/helper"
	"io/ioutil"
	"math"
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
	TESTER_DOMAIN = "default"
	TESTER_APP    = "Tester"
	TESTER_SVC    = "TestService"
	NUM           = 12000
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
	appId, serviceName, version := TESTER_APP, TESTER_SVC, fmt.Sprintf("%d.%d.%d",
		rand.Intn(10), rand.Intn(20), rand.Intn(60))
	r := strings.NewReader(fmt.Sprintf(`{
	"service":{
		"serviceId":"%d",
		"appId":"%s",
		"serviceName":"%s",
		"version":"%s",
		"level":"BACK",
		"status":"UP",
		"schemas":["%s"],
		"properties":{
			"allowCrossApp":"true"
		}
	}
}`, rand.Intn(NUM), appId, serviceName, version, strings.Repeat("x", 160)))

	client := helper.NewClient()
	req, err := http.NewRequest("GET", (&url.URL{
		Scheme:   "http",
		Host:     helper.GetServiceCenterAddress(),
		Path:     EXIST_API,
		RawQuery: "type=microservice&appId=" + appId + "&serviceName=" + serviceName + "&version=" + version,
	}).String(), nil)
	req.Header = http.Header{
		"X-Domain-Name": []string{TESTER_DOMAIN},
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
	resp.Body.Close()

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
		"X-Domain-Name": []string{TESTER_DOMAIN},
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
	resp.Body.Close()
	print(resp.StatusCode != http.StatusOK, time.Now().Sub(t) > time.Second,
		"Create:", string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
}

func RegisterTesterInst() {
	serviceId := fmt.Sprint(rand.Intn(NUM))
	r := strings.NewReader(fmt.Sprintf(`{
	"instance": {
		"serviceId":"%s",
		"endpoints":["rest://127.0.0.1:%d"],
		"hostName":"tester_0_1_2_3",
		"status":"UP",
		"healthCheck":{"mode":"push","interval":30,"times":3},
	}
}`, serviceId, rand.Intn(math.MaxInt16)))
	u := url.URL{
		Scheme: "http",
		Host:   helper.GetServiceCenterAddress(),
		Path:   strings.Replace(REGISTER_API, ":serviceId", serviceId, 1),
	}
	req, err := http.NewRequest("POST", u.String(), r)
	req.Header = http.Header{
		"X-Domain-Name": []string{TESTER_DOMAIN},
		"Content-Type":  []string{"application/json"},
	}

	t := time.Now()
	client := helper.NewClient()
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	print(resp.StatusCode != http.StatusOK, time.Now().Sub(t) > time.Second,
		"Register:", string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
}

func HeartbeatTesterInst() {
	serviceId := fmt.Sprint(rand.Intn(NUM))
	instanceId := fmt.Sprint(rand.Intn(math.MaxInt16))
	u := url.URL{
		Scheme: "http",
		Host:   helper.GetServiceCenterAddress(),
		Path:   strings.Replace(strings.Replace(HEARTBEAT_API, ":serviceId", serviceId, 1), ":instanceId", instanceId, 1),
	}
	req, err := http.NewRequest("PUT", u.String(), nil)
	req.Header = http.Header{
		"X-Domain-Name": []string{TESTER_DOMAIN},
		"Content-Type":  []string{"application/json"},
	}

	t := time.Now()
	client := helper.NewClient()
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	print(resp.StatusCode != http.StatusOK, time.Now().Sub(t) > time.Second,
		"HeartbeatTesterInst:", instanceId, string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
}

func FindTesterInsts() {
	serviceId := fmt.Sprint(NUM)

	versionRules := []string{
		"latest",
		"0+",
		fmt.Sprintf("%d.%d.%d",
			rand.Intn(10), rand.Intn(20), rand.Intn(60)) +
			"-" +
			fmt.Sprintf("%d.%d.%d",
				rand.Intn(10), rand.Intn(20), rand.Intn(60)),
		fmt.Sprintf("%d.%d.%d",
			rand.Intn(10), rand.Intn(20), rand.Intn(60)),
	}
	v := versionRules[rand.Intn(len(versionRules))]

	t := time.Now()
	client := helper.NewClient()
	resp, err := client.Do(&http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme:   "http",
			Host:     helper.GetServiceCenterAddress(),
			Path:     FIND_API,
			RawQuery: "appId=" + TESTER_APP + "&serviceName=" + TESTER_SVC + "&version=" + url.QueryEscape(v),
		},
		Header: http.Header{
			"X-Domain-Name": []string{TESTER_DOMAIN},
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
	resp.Body.Close()
	print(resp.StatusCode != http.StatusOK, time.Now().Sub(t) > time.Second,
		"FindTesterInsts:", v, string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
}

func GetSCInsts() {
	serviceId := helper.GetServiceCenterId()
	t := time.Now()
	client := helper.NewClient()
	resp, err := client.Do(&http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   helper.GetServiceCenterAddress(),
			Path:   strings.Replace(REGISTER_API, ":serviceId", serviceId, 1),
		},
		Header: http.Header{
			"X-Domain-Name": []string{TESTER_DOMAIN},
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
	resp.Body.Close()
	print(resp.StatusCode != http.StatusOK, time.Now().Sub(t) > time.Second,
		"GetSCInsts:", string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
}
