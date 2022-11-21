package common

import (
	"context"
	"net"
)

type RpcInfo struct {
	// 路由名
	ServiceName string
	// interface名
	MethodName string
	// 协议名
	Protocol string
	// 此次rpc的时间
	Timestampe int64
	// 此次rpc的上下文id
	Cid int64
	// error
	Err string
	// Header
	header map[string]string
}

type Request struct {
	// 传入的参数 server端可以不返回，client端如果调用有参函数则需要返回
	Param []interface{}
}

type Response struct {
	// 传出的结果
	Response []interface{}
}

type Kmessage struct {
	rpcInfo *RpcInfo
	in      *Request
	out     *Response
}

// 一次请求的形式 ctx req resp都需要包含
type Endpoint func(ctx context.Context, req *Kmessage) (resp *Kmessage, err error)

// 是krpcServer的高级封装，可以实现不同的krpcServer
type RemoteServer interface {
	Init() error
	DecodeRequest(net.Conn, chan error)
	EncoceResponse(net.Conn, chan error)
	Close() error
}

// 是krpcClient的高级封装，可以实现不同的krpcClient
type RemoteClient interface {
	Init() error
	EncodeRequest(net.Conn, chan error)
	DecodeResponse(net.Conn, chan error)
	Close() error
}
