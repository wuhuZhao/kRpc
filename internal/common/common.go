package common

import (
	"context"
	"kRpc/internal/codec"
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
	Header map[string]string
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
	RpcInfo *RpcInfo
	In      *Request
	Out     *Response
}

// 一次请求的形式 ctx req resp都需要包含
type Endpoint func(ctx context.Context, req *Kmessage) (resp *Kmessage, err error)

// 是krpcServer的高级封装，可以实现不同的krpcServer
type RemoteServer interface {
	Init(co codec.Codec) error
	DecodeRequest(net.Conn, *Kmessage) error
	EncoceResponse(net.Conn, *Kmessage) error
	Close() error
}

// 是krpcClient的高级封装，可以实现不同的krpcClient
type RemoteClient interface {
	Init(co codec.Codec) error
	EncodeRequest(net.Conn, *Kmessage) error
	DecodeResponse(net.Conn, *Kmessage) error
	Close() error
}
