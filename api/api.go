package api

import (
	"fmt"
	"github.com/little-cui/sc-tester/helper"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	HEARTBEAT_API   = "/registry/v3/microservices/:serviceId/instances/:instanceId/heartbeat"
	CREATE_API      = "/registry/v3/microservices"
	REGISTER_API    = "/registry/v3/microservices/:serviceId/instances"
	FIND_API        = "/registry/v3/instances"
	EXIST_API       = "/registry/v3/existence"
	INSTANCE_API    = "/registry/v3/microservices/:serviceId/instances/:instanceId"
	TESTER_DOMAIN   = "default"
	TESTER_APP      = "Tester"
	TESTER_SVC      = "TestService"
	NUM             = 120
	SVC_PER_DOMAIN  = 200
	verPerSvc       = 10
	INST_PER_DOMAIN = 1000
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

func CreateTesterService(d, s int) {
	serviceId := fmt.Sprint(s)
	m := s % (SVC_PER_DOMAIN / verPerSvc)
	v := s % SVC_PER_DOMAIN
	appId, serviceName, version := TESTER_APP, TESTER_SVC, fmt.Sprintf("%s.%d.%d",
		serviceId, m, v)

	schemaIds := ""
	for i := 0; i < 30; i++ {
		schemaIds += fmt.Sprintf(",\"%s\"", strings.Repeat("x", 160))
	}
	r := strings.NewReader(fmt.Sprintf(`{
	"service":{
		"serviceId":"%s",
		"appId":"%s",
		"serviceName":"%s%d",
		"version":"%s",
		"level":"BACK",
		"status":"UP",
		"schemas":[%s],
		"properties":{
			"allowCrossApp":"true"
		}
	}
}`, serviceId, appId, serviceName, m, version, schemaIds[1:]))

	client := helper.NewClient()
	req, err := http.NewRequest("GET", (&url.URL{
		Scheme:   "http",
		Host:     helper.GetServiceCenterAddress(),
		Path:     EXIST_API,
		RawQuery: "type=microservice&appId=" + appId + "&serviceName=" + serviceName + "&version=" + version,
	}).String(), nil)
	req.Header = http.Header{
		"X-Domain-Name": []string{TESTER_DOMAIN + fmt.Sprint(d)},
		"Content-Type":  []string{"application/json"},
	}
	t := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		print(true, time.Now().Sub(t) > time.Second,
			"exist:", appId, serviceName, version, err, "spend:", time.Now().Sub(t))
		return
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
		"X-Domain-Name": []string{TESTER_DOMAIN + fmt.Sprint(d)},
		"Content-Type":  []string{"application/json"},
	}
	t = time.Now()
	resp, err = client.Do(req)
	if err != nil {
		print(true, time.Now().Sub(t) > time.Second,
			"Create:", serviceId, err, "spend:", time.Now().Sub(t))
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	print(resp.StatusCode != http.StatusOK, time.Now().Sub(t) > time.Second,
		"Create:", string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
}

func RegisterTesterInst(i, s, j int) {
	serviceId := fmt.Sprint(s)
	instanceId := fmt.Sprint(j)
	eps := ""
	for i := 0; i < 30; i++ {
		eps += fmt.Sprintf(",\"%s\"", strings.Repeat("x", 160))
	}
	r := strings.NewReader(fmt.Sprintf(`{
	"instance": {
		"serviceId":"%s",
		"instanceId":"%s",
		"endpoints":[%s],
		"hostName":"tester_0_1_2_3",
		"status":"UP",
		"healthCheck":{"mode":"push","interval":30,"times":3}
	}
}`, serviceId, instanceId, eps[1:]))
	u := url.URL{
		Scheme: "http",
		Host:   helper.GetServiceCenterAddress(),
		Path:   strings.Replace(REGISTER_API, ":serviceId", serviceId, 1),
	}
	req, err := http.NewRequest("POST", u.String(), r)
	req.Header = http.Header{
		"X-Domain-Name": []string{TESTER_DOMAIN + fmt.Sprint(i)},
		"Content-Type":  []string{"application/json"},
	}

	t := time.Now()
	client := helper.NewClient()
	resp, err := client.Do(req)
	if err != nil {
		print(true, time.Now().Sub(t) > time.Second,
			"Register:", instanceId, err, "spend:", time.Now().Sub(t))
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	print(resp.StatusCode != http.StatusOK, time.Now().Sub(t) > time.Second,
		"Register:", string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
}

func HeartbeatTesterInst(i, s, j int) {
	serviceId := fmt.Sprint(s)
	instanceId := fmt.Sprint(j)
	u := url.URL{
		Scheme: "http",
		Host:   helper.GetServiceCenterAddress(),
		Path:   strings.Replace(strings.Replace(HEARTBEAT_API, ":serviceId", serviceId, 1), ":instanceId", instanceId, 1),
	}
	req, err := http.NewRequest("PUT", u.String(), nil)
	req.Header = http.Header{
		"X-Domain-Name": []string{TESTER_DOMAIN + fmt.Sprint(i)},
		"Content-Type":  []string{"application/json"},
	}

	t := time.Now()
	client := helper.NewClient()
	resp, err := client.Do(req)
	if err != nil {
		print(true, time.Now().Sub(t) > time.Second,
			"HeartbeatTesterInst:", instanceId, err, "spend:", time.Now().Sub(t))
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	print(resp.StatusCode != http.StatusOK, time.Now().Sub(t) > time.Second,
		"HeartbeatTesterInst:", instanceId, string(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
}

func FindTesterInsts(d, s int, _ int) {
	serviceId := fmt.Sprint(s)
	m := s % (SVC_PER_DOMAIN / verPerSvc)
	appId, serviceName := TESTER_APP, TESTER_SVC+fmt.Sprint(m)
	v := "0+"

	t := time.Now()
	client := helper.NewClient()
	resp, err := client.Do(&http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme:   "http",
			Host:     helper.GetServiceCenterAddress(),
			Path:     FIND_API,
			RawQuery: "appId=" + appId + "&serviceName=" + serviceName + "&version=" + url.QueryEscape(v),
		},
		Header: http.Header{
			"X-Domain-Name": []string{TESTER_DOMAIN + fmt.Sprint(d)},
			"X-ConsumerId":  []string{serviceId},
		},
	})
	if err != nil {
		print(true, time.Now().Sub(t) > time.Second,
			"FindTesterInsts:", appId, serviceName, v, err, "spend:", time.Now().Sub(t))
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(true, time.Now().Sub(t) > time.Second,
			"FindTesterInsts:", appId, serviceName, v, "read body:", err, "spend:", time.Now().Sub(t))
		return
	}
	resp.Body.Close()
	print(resp.StatusCode != http.StatusOK, time.Now().Sub(t) > time.Second,
		"FindTesterInsts:", appId, serviceName, v, serviceId, "resp:", len(body), "status:", resp.StatusCode, "spend:", time.Now().Sub(t))
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
