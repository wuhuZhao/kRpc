package core

import (
	"context"
	"net"
)

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
	Call(conn net.Conn, ep Endpoint, request *Message) (*Message, error)
}
