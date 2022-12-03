package discovry

import (
	"fmt"
	"kRpc/pkg/klog"
	"math/rand"

	consulapi "github.com/hashicorp/consul/api"
)

type DiscoveryConfig struct {
	ID      string
	Name    string
	Tags    []string
	Port    int
	Address string
	UseTCP  bool
}

var consulAddress = "127.0.0.1:8500"

func SetConsulAddress(address string) {
	consulAddress = address
}

func RegisterService(dis DiscoveryConfig) error {
	config := consulapi.DefaultConfig()
	config.Address = consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		klog.Errf("create consul client : %v\n", err.Error())
	}
	registration := &consulapi.AgentServiceRegistration{
		ID:      dis.ID,
		Name:    dis.Name,
		Port:    dis.Port,
		Tags:    dis.Tags,
		Address: dis.Address,
	}
	if dis.UseTCP {
		check := &consulapi.AgentServiceCheck{}
		check.TCP = fmt.Sprintf("%s:%d", registration.Address, registration.Port)
		check.Timeout = "5s"
		check.Interval = "5s"
		check.DeregisterCriticalServiceAfter = "60s"
		registration.Check = check
	}

	if err := client.Agent().ServiceRegister(registration); err != nil {
		klog.Errf("register to consul error: %v\n", err.Error())
		return err

	}
	return nil
}

func Discovery(serviceName string) []*consulapi.ServiceEntry {
	config := consulapi.DefaultConfig()
	config.Address = consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		klog.Errf("consul client error: %v", err)
	}
	service, _, err := client.Health().Service(serviceName, "", false, nil)
	if err != nil {
		klog.Errf("consul client get serviceIp error: %v", err)
	}
	return service
}

func SelectOnePod(psm string) *consulapi.ServiceEntry {
	total := Discovery(psm)
	return total[rand.Intn(len(total))]
}
