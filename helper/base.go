package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
)

const (
	SERVER_ADDRESS        = "127.0.0.1:30100"
	REGISTRY_APP_ID       = "default"
	REGISTRY_SERVICE_NAME = "SERVICECENTER"
	REGISTRY_VERSION      = "3.0.0"
)

type ServiceExistResponse struct {
	ServiceId string `json:"serviceId"`
}

type InstancesResponse struct {
	Instances []*Instance `json:"instances,omitempty"`
}

type Instance struct {
	InstanceId string `json:"instanceId"`
}

func GetServiceCenterAddress() string {
	return SERVER_ADDRESS
}

func GetServiceCenterId() string {
	u := url.URL{
		Scheme: "http",
		Host:   SERVER_ADDRESS,
		Path:   "/registry/v3/existence",
		RawQuery: "type=microservice&appId=" + REGISTRY_APP_ID +
			"&serviceName=" + REGISTRY_SERVICE_NAME +
			"&version=" + REGISTRY_VERSION,
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err)
	}
	req.Header = http.Header{
		"X-Domain-Name": []string{"default"},
		"Content-Type":  []string{"application/json"},
	}
	resp, err := http.DefaultClient.Do(req)
	if resp == nil {
		fmt.Println("GetServiceCenterId:", resp, err)
		return ""
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		// panic(string(respBody))
		fmt.Println("GetServiceCenterId:", resp.StatusCode, err)
		return ""
	}
	var serviceResponse ServiceExistResponse
	err = json.Unmarshal(respBody, &serviceResponse)
	if err != nil {
		panic(err)
	}
	return serviceResponse.ServiceId
}

func GetServiceCenterInstanceId(serviceId string) string {
	u := url.URL{
		Scheme: "http",
		Host:   SERVER_ADDRESS,
		Path:   "/registry/v3/microservices/" + serviceId + "/instances",
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err)
	}
	req.Header = http.Header{
		"X-Domain-Name": []string{"default"},
		"Content-Type":  []string{"application/json"},
	}
	req.Header.Set("X-ConsumerId", serviceId)
	resp, err := http.DefaultClient.Do(req)
	if resp == nil {
		fmt.Println("GetServiceCenterInstanceId:", resp, err)
		return ""
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("GetServiceCenterInstanceId:", resp.StatusCode, err)
		return ""
	}
	var instancesResponse InstancesResponse
	err = json.Unmarshal(respBody, &instancesResponse)
	if err != nil {
		panic(err)
	}
	l := len(instancesResponse.Instances)
	if l == 0 {
		return ""
	}
	return instancesResponse.Instances[rand.Intn(l)].InstanceId
}
