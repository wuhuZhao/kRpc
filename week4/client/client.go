package client

import (
	"context"
	"net"
	"week4/core"
)

type Option struct {
	serverIp       string
	serverPort     string
	serverProtocol string
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
	resp, err := c.remoteClient.Call(c.conn, c.eps, request)
	if err != nil {
		return err
	}
	response = resp
	return nil
}

// 将mds做成递归的链式调用，与server的原理是一致的，但是一开始的传参不是invokeHandler, invokeHandler应该是最后一步
func (c *Client) chain(next core.Endpoint) core.Endpoint {
	for i := 1; i < len(c.mds); i++ {
		next = c.mds[i](next)
	}
	return next
}

// client端真正的调用
func (c *Client) Call(request *core.Message) (response *core.Message, err error) {
	// todo 将中间件构建成一个递归调用
	if c.eps == nil {
		// 将实现的InvokeHandler放进去当做最后一步的封装，就是优先执行
		c.mds = append(c.mds, func(e core.Endpoint) core.Endpoint {
			return func(ctx context.Context, req, resp *core.Message) error {
				err := c.invoke(ctx, req, resp)
				if err != nil {
					return err
				}
				return e(ctx, req, resp)
			}
		})
		// len == 1时，只有invokehandler,所以不需要转成链
		if len(c.mds) == 1 {
			c.eps = c.invoke
		} else {
			// 注册一下虚假的中间件去让first能包裹住，但是这个不影响结果，类似上面的e(ctx,req,resp),然后即封装出一个调用链去实现从remoteCall -> middlewares的调用
			c.eps = c.chain(c.mds[0](func(ctx context.Context, req *core.Message, resp *core.Message) error {
				return nil
			}))
		}
	}
	err = c.eps(context.Background(), request, response)
	return
}

// 中间件插入
func (c *Client) Use(md core.Middleware) {
	c.mds = append(c.mds, md)
}

func NewDefaultServer(option *Option) (*Client, error) {
	client := &Client{}
	conn, err := net.Dial(option.serverProtocol, option.serverIp+":"+option.serverPort)
	if err != nil {
		return nil, err
	}
	remoteClient := core.NewKrpcClient()
	client.init(remoteClient, conn)
	return client, nil
}
