package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"reflect"
	"sync"
	"time"
	"week5/core"
	"week5/internal/klog"

	"gopkg.in/yaml.v2"
)

type yamlConfig struct {
	Service struct {
		Psm      string `yaml: "psm"`
		Port     string `yaml: "port"`
		Ip       string `yaml: "ip"`
		Protocol string `yaml: "protocol"`
	} `yaml: "service"`
}

type option struct {
	ServerIp       string
	ServerPort     string
	ServerProtocol string
	// 用于自定义的customizeProps 方便开发者自己自定义启动Hooks
	CustomizeProps map[string]interface{}
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
	// 启动时调用各种组件，仅在启动时调用
	shs []StartHook
	m   *sync.Mutex
	// 保存相关配置
	opt *option
	// meta yaml里存在的文件
	meta []byte
}

// 真正初始化的地方,通用化ln和remoteServer, Option里的配置也保存一下，用于钩子函数的实现也可以注册相关自定义启动的钩子函数
func (s *Server) init(ln net.Listener, remoteServer core.RemoteServer, opt *option) {
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

// 添加中间件, 不断dfs下去，func套func完成中间件的调用层级，init的时候为reflect调用的ep, 用slice append
func (s *Server) Use(mdw core.Middleware) {
	s.mds = append(s.mds, mdw)
}

//   不需要修改  request -> middleware-> send
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
		// 调用钩子函数注册
		for i := 0; i < len(s.shs); i++ {
			s.shs[i](s.opt)
		}
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

func (s *Server) getYamlOption(filePath string) *option {
	opt := &option{}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		klog.Errf("[krpc]get yaml config err: %v", err.Error())
		return s.getDefaultConfig()
	}
	var config yamlConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		klog.Errf("[krpc]get yaml config err: %v", err.Error())
		return s.getDefaultConfig()
	}
	opt.ServerIp = config.Service.Ip
	opt.ServerPort = config.Service.Port
	opt.ServerProtocol = config.Service.Protocol
	createRegisterOption(opt, &RegistrationConfig{
		Name: config.Service.Psm,
		ID:   time.Now().GoString(),
	})
	klog.Infof("%+v\n", opt)
	return opt
}

func (s *Server) getDefaultConfig() *option {
	opt := &option{}
	opt.ServerIp = "172.18.0.1"
	opt.ServerPort = "10011"
	opt.ServerProtocol = "tcp"
	createRegisterOption(s.opt, &RegistrationConfig{
		Name: "defaultServer",
		ID:   time.Now().GoString(),
	})
	klog.Infof("[krpc] there is not valid filePath in project, use default config! default config: %v", s.opt)
	return opt
}

// 返回默认的服务端，使用krpcServer作为底层去使用
func NewDefaultServer(configPath string) (*Server, error) {
	krpcServer, err := core.NewKrpcServer()
	if err != nil {
		return nil, err
	}
	server := &Server{}
	server.opt = server.getYamlOption(configPath)
	listen, err := net.Listen(server.opt.ServerProtocol, server.opt.ServerIp+":"+server.opt.ServerPort)
	if err != nil {
		return nil, err
	}
	server.init(listen, krpcServer, server.opt)
	// 注册服务到consul里
	server.shs = append(server.shs, RegisterService)
	return server, nil
}
