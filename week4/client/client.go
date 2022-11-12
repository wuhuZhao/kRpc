package client

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"week4/core"

	consulapi "github.com/hashicorp/consul/api"
)

const consulAddress = "127.0.0.1:8500"

type Option struct {
	ServerIp       string
	ServerPort     string
	ServerProtocol string
	ServiceName    string
}

type Client struct {
	remoteClient  core.RemoteClient
	conn          net.Conn
	eps           core.Endpoint
	mds           []core.Middleware
	invokeHandler core.Endpoint
}

func (c *Client) init(rc core.RemoteClient, conn net.Conn) {
	c.remoteClient = rc
	c.conn = conn
	c.invokeHandler = c.invoke
	c.mds = []core.Middleware{}
}

// 做高级封装的call，实际底层还是调用remoteClient的call方法，把他封装到Middleware中方便使用
func (c *Client) invoke(ctx context.Context, request *core.Message, response *core.Message) error {
	err := c.remoteClient.Call(c.conn, request, response)
	if err != nil {
		return err
	}
	return nil
}

//  client应该是request -> Middleware->send  我搞反了 这个需要修改
func (c *Client) chain(next core.Endpoint) core.Endpoint {
	for i := len(c.mds) - 1; i >= 0; i-- {
		next = c.mds[i](next)
	}
	return next
}

// client端真正的调用
func (c *Client) Call(request *core.Message) (*core.Message, error) {
	// todo 将中间件构建成一个递归调用
	if c.eps == nil {
		c.eps = c.chain(c.invokeHandler)
	}
	response := &core.Message{RpcInfo: &core.RpcInfo{}}
	err := c.eps(context.Background(), request, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 中间件插入
func (c *Client) Use(md core.Middleware) {
	c.mds = append(c.mds, md)
}

// consul服务发现
func (c *Client) Discovery(serviceName string) *consulapi.ServiceEntry {
	config := consulapi.DefaultConfig()
	config.Address = consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("consul client error : ", err)
	}

	service, _, err := client.Health().Service(serviceName, "", false, nil)
	if err != nil {
		fmt.Println("consul client get serviceIp error : ", err)
	}
	return service[0]
}

func NewDefaultClient(option *Option) (*Client, error) {
	client := &Client{}
	if option.ServiceName != "" {
		service := client.Discovery(option.ServiceName)
		option.ServerIp = service.Service.Address
		option.ServerPort = strconv.Itoa(service.Service.Port)
	}
	conn, err := net.Dial(option.ServerProtocol, option.ServerIp+":"+option.ServerPort)
	if err != nil {
		return nil, err
	}
	remoteClient := core.NewKrpcClient()
	client.init(remoteClient, conn)
	return client, nil
}
