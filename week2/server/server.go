package server

import (
	"context"
	"fmt"
	"net"
	"time"
	"week2/codec"
)

type EndPoint func(ctx context.Context, req interface{}, resp interface{}) error

type option struct {
	connectTime      time.Duration
	mux              bool
	bufferSize       int
	registerProtocol string
	registerDomain   string
	defaultMode      bool
	serverIp         string
	serverPort       string
	serverProtocol   string
}

func NewKrpcOptionWithDefaultMode(serverIp, serverPort, serverProtocol string) *option {
	return &option{defaultMode: true, serverIp: serverIp, serverPort: serverPort, serverProtocol: serverProtocol}
}

type krpcServer struct {
	opt *option
	ln  net.Listener
	eps EndPoint
	co  codec.Codec
}

func (k *krpcServer) init() error {
	listen, err := net.Listen(k.opt.serverProtocol, k.opt.serverIp+":"+k.opt.serverPort)
	if err != nil {
		return err
	}
	k.ln = listen
	return nil
}

func (k *krpcServer) Serve() chan error {
	err := make(chan error, 100)
	var handler func(conn net.Conn, ch chan error)
	handler = func(conn net.Conn, ch chan error) {
		var resp interface{}
		if er := k.co.Decode(conn, &resp); er != nil {
			ch <- er
		}
		fmt.Printf("[kRpcServer] get resp %v \n", resp)
		// todo 实现调用哪个方法
		req := map[string]string{
			"context": "hello, I am Krpc",
			"date":    time.UnixDate,
		}
		if er := k.co.Encode(conn, req); er != nil {
			ch <- er
		}
		fmt.Printf("[kRpcServer] reply req %v \n", req)
	}
	go func(ln net.Listener, c chan error) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				c <- err
				continue
			}
			go handler(conn, c)

		}
	}(k.ln, err)
	return err
}

func NewKrpcServer(opt *option) *krpcServer {
	server := &krpcServer{}
	server.opt = opt
	// 初始化
	err := server.init()
	if err != nil {
		panic(fmt.Sprintf("[kRpc] the server can't start! %+v", err))
	}
	// 走目前的逻辑
	if server.opt.defaultMode {
		server.co = codec.NewJsonCodec(&codec.Option{})
		//todo 需要走consul的注册逻辑
	} else {
		//todo 使用其他注册中心的逻辑
	}
	return server

}
