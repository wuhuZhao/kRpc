package server

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"sync"
	"week4/core"
)

type Option struct {
	serverIp       string
	serverPort     string
	serverProtocol string
}

type Server struct {
	remoteServer core.RemoteServer
	// service 对应Idl的service, 同一个service下面有不同的interface
	service map[string]map[string]struct{}
	// 一个service+ interface 对应一个真正的函数调用
	method        map[string]reflect.Value
	invokeHandler core.Endpoint
	eps           core.Endpoint
	mds           []core.Middleware
	// 可以选择自己的net库，但目前不开放，转向使用Option去直接给生成
	ln net.Listener
	m  *sync.Mutex
}

// 真正初始化的地方,通用化ln和remoteServer
func (s *Server) init(ln net.Listener, remoteServer core.RemoteServer) {
	s.invokeHandler = s.serverEndpoint
	s.ln = ln
	s.remoteServer = remoteServer
	s.service = map[string]map[string]struct{}{}
	s.method = map[string]reflect.Value{}
	s.mds = []core.Middleware{}
	s.m = &sync.Mutex{}

}

// 调用真正反射的方法去完成一次rpc调用，不放在中间件中，但是以中间件的形式
func (s *Server) serverEndpoint(ctx context.Context, request *core.Message, response *core.Message) error {
	if service, ex := s.service[request.ServiceName]; ex {
		if _, ok := service[request.MethodName]; ok {
			if fn, ok := s.method[request.ServiceName+request.MethodName]; ok {
				paramValue := []reflect.Value{}
				for i := 0; i < len(request.Param); i++ {
					paramValue = append(paramValue, reflect.ValueOf(request.Param[i]))
				}
				out := fn.Call(paramValue)
				resValue := []interface{}{}
				for i := 0; i < len(out); i++ {
					resValue = append(resValue, out[i].Interface())
				}
				response.Response = resValue
			} else {
				response.Err = fmt.Sprintf("[KrpcServer] can't find the method %s in method table", request.MethodName)
			}
		} else {
			response.Err = fmt.Sprintf("[KrpcServer] can't find the method %s in service table", request.MethodName)
		}
	} else {
		response.Err = fmt.Sprintf("[KrpcServer] can't find the service %s in service table", request.ServiceName)
	}
	return nil
}

// 添加中间件, 不断dfs下去，func套func完成中间件的调用层级，init的时候为reflect调用的ep, 用slice append 然后通过for循环去组装dfs 不要用递归，这样就能控制递归的顺序了，client和server应该相反
// func (s *Server) Use(next core.Middleware) {
// 	s.ep = next(s.ep)
// }

// 添加中间件, 不断dfs下去，func套func完成中间件的调用层级，init的时候为reflect调用的ep, 用slice append
func (s *Server) Use(mdw core.Middleware) {
	s.mds = append(s.mds, mdw)
}

// 然后通过for循环去组装dfs 不要用递归，这样就能控制递归的顺序了，client和server应该相反
func (s *Server) chain(ep core.Endpoint) core.Endpoint {
	for i := len(s.mds) - 1; i >= 0; i-- {
		ep = s.mds[i](ep)
	}
	return ep
}

func (s *Server) Serve() chan error {
	errChannel := make(chan error, 20)
	// 初始化调用链eps
	s.eps = s.chain(s.invokeHandler)
	// 经典阻塞式IO模型，这样在gorutine多的时候会导致很多的调度开销
	go func(ln net.Listener, c chan error) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				c <- err
				continue
			}
			// 这里调用真正的transport transport需要接受调用中间件和反射调用最终输出resp
			go s.remoteServer.Transport(conn, s.eps, c)
		}
	}(s.ln, errChannel)
	return errChannel
}

// 注册实现类
func (s *Server) Register(serviceImpl interface{}) {
	s.m.Lock()
	defer func() {
		s.m.Unlock()
	}()
	// todo 后续优化一下，可以弄成只用一个Ptr或者直接一个type去实现，否则太乱了
	serviceValue := reflect.ValueOf(serviceImpl)
	serviceValuePtr := serviceValue.Elem()
	serviceName := serviceValuePtr.Type().Name()
	s.service[serviceName] = make(map[string]struct{})
	for i := 0; i < serviceValue.NumMethod(); i++ {
		// Method通过Type去拿才能拿到的
		methodName := serviceValue.Type().Method(i).Name
		s.service[serviceName][methodName] = struct{}{}
		s.method[serviceName+methodName] = serviceValue.Method(i)
	}
}

// 返回默认的服务端，使用krpcServer作为底层去使用
func NewDefaultServer(opt *Option) (*Server, error) {
	listen, err := net.Listen(opt.serverProtocol, opt.serverIp+":"+opt.serverPort)
	if err != nil {
		return nil, err
	}
	krpcServer, err := core.NewKrpcServer()
	if err != nil {
		return nil, err
	}
	server := &Server{}
	server.init(listen, krpcServer)
	return server, nil
}
