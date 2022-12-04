package server

import (
	discovry "kRpc/pkg/discovery"
	"kRpc/pkg/klog"
)

type ServerHooks interface {
	Start()
}

type RegiserHooks struct {
	SerivceName   string
	ServiceIp     string
	ServicePort   int
	ConsulAddress string
	UseTcp        bool
}

func (r *RegiserHooks) Start() {
	if len(r.ConsulAddress) != 0 {
		discovry.SetConsulAddress(r.ConsulAddress)
	}
	if err := discovry.RegisterService(discovry.DiscoveryConfig{
		Name:    r.SerivceName,
		Address: r.ServiceIp,
		Port:    r.ServicePort,
		UseTCP:  r.UseTcp,
	}); err != nil {
		klog.Errf("register to consul error: %v\n", err.Error())
	}
}
