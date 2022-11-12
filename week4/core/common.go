package core

import (
	"context"
	"net"
)

// 一次请求的形式 ctx req resp都需要包含
type Endpoint func(ctx context.Context, req *Message, resp *Message) error

// 是krpcServer的高级封装，可以实现不同的krpcServer
type RemoteServer interface {
	Init() error
	Transport(net.Conn, Endpoint, chan error)
	// Stop方法暂时不实现
	// Stop()
}

// 是krpcClient的高级封装，可以实现不同的krpcClient
type RemoteClient interface {
	Call(conn net.Conn, request *Message, resp *Message) error
}
