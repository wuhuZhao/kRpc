package main

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
)

func Discovery(serviceName string) []*consulapi.ServiceEntry {
	config := consulapi.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Printf("consul client error: %v", err)
	}
	service, _, err := client.Health().Service(serviceName, "", false, nil)
	if err != nil {
		fmt.Printf("consul client get serviceIp error: %v", err)
	}
	return service
}

func TestDiscoeryFromConsul(t *testing.T) {
	t.Logf("client discovery start")
	se := Discovery("main_service")
	for i := 0; i < len(se); i++ {
		t.Logf("the instance Node is %+v\n", se[i].Node)
		t.Logf("the isntance Service is %+v\n", se[i].Service)
		t.Logf("\n")
	}
}
