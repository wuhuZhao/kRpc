package server

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"sync"
	"time"
	"week3/codec"
	"week3/core"
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
	// service 对应Idl的service, 同一个service下面有不同的interface
	service map[string]map[string]struct{}
	// 一个service+ interface 对应一个真正的函数调用
	method map[string]reflect.Value
	mutex  *sync.Mutex
}

func (k *krpcServer) init() error {
	listen, err := net.Listen(k.opt.serverProtocol, k.opt.serverIp+":"+k.opt.serverPort)
	if err != nil {
		return err
	}
	k.ln = listen
	k.mutex = &sync.Mutex{}
	k.service = make(map[string]map[string]struct{})
	k.method = make(map[string]reflect.Value)
	return nil
}

// 暂时使用public方法进行显式调用
func (k *krpcServer) Register(serviceImpl interface{}) {
	k.mutex.Lock()
	defer func() {
		k.mutex.Unlock()
	}()
	serviceValue := reflect.ValueOf(serviceImpl)
	serviceValuePtr := serviceValue.Elem()
	serviceName := serviceValuePtr.Type().Name()
	k.service[serviceName] = make(map[string]struct{})
	for i := 0; i < serviceValue.NumMethod(); i++ {
		// Method通过Type去拿才能拿到的
		methodName := serviceValue.Type().Method(i).Name
		k.service[serviceName][methodName] = struct{}{}
		k.method[serviceName+methodName] = serviceValue.Method(i)
	}
}

func (k *krpcServer) Serve() chan error {
	err := make(chan error, 100)
	var handler func(conn net.Conn, ch chan error)
	handler = func(conn net.Conn, ch chan error) {
		var req core.Message
		if er := k.co.Decode(conn, &req); er != nil {
			ch <- er
		}
		fmt.Printf("[kRpcServer] get req %v \n", req)
		request := req
		// 判断方法是否存在，再选择调用
		if service, ex := k.service[request.ServiceName]; ex {
			if _, ok := service[request.MethodName]; ok {
				if fn, ok := k.method[request.ServiceName+request.MethodName]; ok {
					paramValue := []reflect.Value{}
					for i := 0; i < len(request.Param); i++ {
						paramValue = append(paramValue, reflect.ValueOf(request.Param[i]))
					}
					out := fn.Call(paramValue)
					resp := core.NewWithSuccessMessage(&core.RpcInfo{ServiceName: request.ServiceName, MethodName: request.MethodName, Protocol: request.Protocol})
					res := []interface{}{}
					for i := 0; i < len(out); i++ {
						res = append(res, out[i].Interface())
					}
					resp.Response = res
					if er := k.co.Encode(conn, resp); er != nil {
						ch <- er
					}

				} else {
					resp := core.NewWithFailMessage(&core.RpcInfo{ServiceName: request.ServiceName, MethodName: request.MethodName, Protocol: request.Protocol})
					resp.Err = fmt.Sprintf("[KrpcServer] can't find the method %s in method table", request.MethodName)
					fmt.Printf("%s\n", resp.Err)
					if er := k.co.Encode(conn, resp); er != nil {
						ch <- er
					}
				}
			} else {
				resp := core.NewWithFailMessage(&core.RpcInfo{ServiceName: request.ServiceName, MethodName: request.MethodName, Protocol: request.Protocol})
				resp.Err = fmt.Sprintf("[KrpcServer] can't find the method %s in service table", request.MethodName)
				fmt.Printf("%s\n", resp.Err)
				if er := k.co.Encode(conn, resp); er != nil {
					ch <- er
				}
			}
		} else {
			resp := core.NewWithFailMessage(&core.RpcInfo{ServiceName: request.ServiceName, MethodName: request.MethodName, Protocol: request.Protocol})
			resp.Err = fmt.Sprintf("[KrpcServer] can't find the service %s in service table", request.ServiceName)
			fmt.Printf("%s\n", resp.Err)
			if er := k.co.Encode(conn, resp); er != nil {
				ch <- er
			}
		}
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
