package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

const (
	SERVER_ADDRESS        = "127.0.0.1:30100"
	REGISTRY_APP_ID       = "default"
	REGISTRY_SERVICE_NAME = "SERVICECENTER"
	REGISTRY_VERSION      = "0.0.1"
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
			"&version=latest&env=development",
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err)
	}
	req.Header = http.Header{
		"X-Domain-Name": []string{"default"},
		"Content-Type":  []string{"application/json"},
	}
	client := NewClient()
	resp, err := client.Do(req)
	if resp == nil {
		fmt.Println("GetServiceCenterId:", resp, err)
		return ""
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// panic(string(respBody))
		fmt.Println("GetServiceCenterId:", resp.StatusCode, string(respBody), err)
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
	client := NewClient()
	resp, err := client.Do(req)
	if resp == nil {
		fmt.Println("GetServiceCenterInstanceId:", resp, err)
		return ""
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
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

var c *http.Client
var once sync.Once

func NewClient() *http.Client {
	once.Do(func() {
		c = &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        100000,
				MaxIdleConnsPerHost: 100000,
			},
		}
	})
	return c
}
