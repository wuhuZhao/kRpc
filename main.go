package main

import (
	"kRpc/pkg/klog"
	"kRpc/server"
)

type MyImple struct{}

func (m *MyImple) Add(a, b int) int {
	return a + b
}

func main() {
	s, err := server.NewDefaultServer(server.CreateOption(8999, "192.168.0.124", "tcp", "haozhao.test"))
	if err != nil {
		klog.Errf("start krpc error: %v\n", err.Error())
		return
	}
	s.Register(&MyImple{})
	s.Start()
}
