package client

import (
	"context"
	"kRpc/internal/common"
	"kRpc/internal/core"
	"kRpc/internal/middleware"
	discovry "kRpc/pkg/discovery"
	"net"
	"sync"
)

type Client struct {
	remoteClient common.RemoteClient
	psm          string
	// 基于consul的话 一个服务可能会有多个ip  用ip区分conn
	conn           map[string]*net.Conn
	invoke         common.Endpoint
	mds            []middleware.Middleware
	mutex          *sync.Mutex
	invokeComplete bool
}

func (c *Client) Use(md middleware.Middleware) {
	c.mds = append(c.mds, md)
}

// client端的调用
func (c *Client) handler(ctx context.Context, req *common.Kmessage) (*common.Kmessage, error) {
	entry := discovry.SelectOnePod(c.psm)
	//todo 考虑用sync.Map 因为不安全
	if _, ok := c.conn[entry.Node.Address]; !ok {
		conn, err := net.Dial("tcp", entry.Node.Address)
		if err != nil {
			return nil, err
		}
		c.conn[entry.Node.Address] = &conn
	}
	conn := c.conn[entry.Node.Address]
	if err := c.remoteClient.EncodeRequest(*conn, req); err != nil {
		return nil, err
	}
	resp := &common.Kmessage{}
	if err := c.remoteClient.DecodeResponse(*conn, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// 构造一个中间件的构造逻辑
func (c *Client) Call(req *common.Kmessage) (*common.Kmessage, error) {
	if !c.invokeComplete {
		c.mutex.Lock()
		defer func() { c.mutex.Unlock() }()
		if !c.invokeComplete {
			c.invoke = middleware.WrapperMiddleware(c.mds, c.handler)
			c.invokeComplete = true
		}
	}
	return c.invoke(context.Background(), req)
}

// 创建一个需要链接的客户端
func NewDefaultClient(psm string) (*Client, error) {
	client := &Client{psm: psm, remoteClient: core.NewDefaultKrpcClient()}
	return client, nil
}
