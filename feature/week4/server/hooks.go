package server

import (
	"fmt"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
)

// consulAddress的默认地址
const consulAddress = "127.0.0.1:8500"

// 启动的钩子函数
type StartHook func(opt *Option)

// 结构体，方便opt的构建
type RegistrationConfig struct {
	ID   string
	Name string
	Tags []string
}

// 创建consul的option
func createRegisterOption(opt *Option, reg *RegistrationConfig) {
	if opt.CustomizeProps == nil {
		opt.CustomizeProps = map[string]interface{}{}
	}
	opt.CustomizeProps["ID"] = reg.ID
	opt.CustomizeProps["Name"] = reg.Name
	opt.CustomizeProps["Tags"] = reg.Tags
}

// 注册consul的starthook
func RegisterService(opt *Option) {
	config := consulapi.DefaultConfig()
	config.Address = consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("consul client error : ", err)
	}
	registration := &consulapi.AgentServiceRegistration{}
	registration.ID = opt.CustomizeProps["ID"].(string)
	registration.Name = opt.CustomizeProps["Name"].(string)
	registration.Port, _ = strconv.Atoi(opt.ServerPort)
	registration.Tags = opt.CustomizeProps["Tags"].([]string)
	registration.Address = opt.ServerIp

	// 添加consul健康检查回调函数 目前服务注册暂时不能识别我的包，后续研究一下，给consul回包  证明keeplive
	check := &consulapi.AgentServiceCheck{}
	check.HTTP = fmt.Sprintf("http://%s:%s", registration.Address, "10022")
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "60s"
	registration.Check = check

	if err := client.Agent().ServiceRegister(registration); err != nil {
		fmt.Printf("register to consul error: %v \n", err)
	}
}
