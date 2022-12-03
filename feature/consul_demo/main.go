package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"

	consulapi "github.com/hashicorp/consul/api"
)

type DiscoveryConfig struct {
	ID      string
	Name    string
	Tags    []string
	Port    int
	Address string
}

var consulAddress = "127.0.0.1:8500"

func RegisterService(dis DiscoveryConfig) error {
	config := consulapi.DefaultConfig()
	config.Address = consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Printf("create consul client : %v\n", err.Error())
	}
	registration := &consulapi.AgentServiceRegistration{
		ID:      dis.ID,
		Name:    dis.Name,
		Port:    dis.Port,
		Tags:    dis.Tags,
		Address: dis.Address,
	}
	// 启动tcp的健康检测，注意address不能使用127.0.0.1或者localhost，因为consul-agent在docker容器里，如果用这个的话，
	// consul会访问容器里的port就会出错，一直检查不到实例
	check := &consulapi.AgentServiceCheck{}
	check.HTTP = fmt.Sprintf("http://%s:%d/", registration.Address, registration.Port)
	// check.TCP = fmt.Sprintf("%s:%d", registration.Address, registration.Port)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "60s"
	registration.Check = check

	if err := client.Agent().ServiceRegister(registration); err != nil {
		fmt.Printf("register to consul error: %v\n", err.Error())
		return err
	}
	return nil
}

func startHttp() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("consul get uri: %s\n", r.RequestURI)
		w.Write([]byte("hello consul"))
	})
	if err := http.ListenAndServe(":10111", nil); err != nil {
		fmt.Printf("start http server error: %v\n", err)
	}
}

func startTcp() {
	ls, err := net.Listen("tcp", ":10111")
	if err != nil {
		fmt.Printf("start tcp listener error: %v\n", err.Error())
		return
	}
	for {
		conn, err := ls.Accept()
		if err != nil {
			fmt.Printf("connect error: %v\n", err.Error())
		}
		go func(conn net.Conn) {
			_, err := bufio.NewWriter(conn).WriteString("hello consul")
			if err != nil {
				fmt.Printf("write conn error: %v\n", err)
			}
		}(conn)
	}
}
func main() {
	ch := make(chan error)
	dis := DiscoveryConfig{
		ID:      "9527",
		Name:    "main_service",
		Tags:    []string{"a", "b"},
		Port:    10111,
		Address: "192.168.0.124", //通过ifconfig查看本机的eth0的ipv4地址
	}
	// go startTcp()
	go startHttp()
	RegisterService(dis)
	// 阻塞等待
	<-ch
}
