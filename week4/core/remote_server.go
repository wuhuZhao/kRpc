package core

import (
	"context"
	"errors"
	"fmt"
	"net"
	"week4/codec"
)

var _ RemoteServer = (*krpcServer)(nil)

type krpcServer struct {
	co codec.Codec
}

// init方法 目前只加载encode
func (k *krpcServer) Init() error {
	k.co = codec.NewJsonCodec(&codec.Option{})
	return nil
}

// 真正的transport调用，只需要传入封装好的中间件以及call，最后实现一个 encode-> (middleware -> call) - > decode 的模型
func (k *krpcServer) Transport(conn net.Conn, ep Endpoint, ch chan error) {
	var req *Message
	var resp *Message
	if er := k.co.Decode(conn, req); er != nil {
		ch <- er
	}
	if err := ep(context.Background(), req, resp); err != nil {
		ch <- err
	}
	if er := k.co.Encode(conn, resp); er != nil {
		ch <- er
	}
}

func NewKrpcServer() (*krpcServer, error) {
	server := &krpcServer{}
	// 初始化
	err := server.Init()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[kRpc] the server can't start! %+v", err))
	}
	return server, nil

}
