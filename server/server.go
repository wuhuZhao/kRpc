package server

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"kRpc/internal/common"
	"kRpc/internal/core"
	"kRpc/internal/middleware"
	"kRpc/pkg/klog"
	"net"
	"reflect"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

type yamlConfig struct {
	Service struct {
		Psm      string `yaml:"psm"`
		Port     int    `yaml:"port"`
		Ip       string `yaml:"ip"`
		Protocol string `yaml:"protocol"`
	} `yaml:"service"`
	Discovery struct {
		Address string `yaml:"address"`
	} `yaml:"discovery"`
}

type option struct {
	ServerIp         string
	ServerPort       int
	ServerProtocol   string
	DiscoveryAddress string
	psm              string
	// 用于自定义的customizeProps 方便开发者自己自定义启动Hooks
	CustomizeProps map[string]interface{}
}

type Server struct {
	remoteServer common.RemoteServer
	method       map[string]map[string]reflect.Value
	invoke       common.Endpoint
	mds          []middleware.Middleware
	ls           *net.Listener
	// meta 暂时为使用，应该是配置中心等信息存储在meta当中
	meta []byte
	// lock
	mutex      *sync.Mutex
	startHooks []ServerHooks
	endHooks   []ServerHooks
}

// 注入ln和使用的remoteServer实现
func (s *Server) init(ln *net.Listener, remoteServer common.RemoteServer) {
	s.ls = ln
	s.remoteServer = remoteServer
	s.invoke = s.defaultInvoke
	//todo start hook
	for i := 0; i < len(s.startHooks); i++ {
		go s.startHooks[i].Start()
	}
}

// 注册中间件
func (s *Server) Use(md middleware.Middleware) {
	s.mds = append(s.mds, md)
}

// 启动hooks
func (s *Server) AddStartHooks(cur ServerHooks) {
	s.startHooks = append(s.startHooks, cur)
}

// 结束hooks
func (s *Server) AddEndHooks(cur ServerHooks) {
	s.endHooks = append(s.endHooks, cur)
}

// 注册rpc方法
func (s *Server) Register(stub interface{}) {
	s.mutex.Lock()
	defer func() { s.mutex.Unlock() }()
	stubValue := reflect.ValueOf(stub)
	stubName := stubValue.Elem().Type().Name()
	s.method[stubName] = map[string]reflect.Value{}
	for i := 0; i < stubValue.NumMethod(); i++ {
		methodName := stubValue.Type().Method(i).Name
		s.method[stubName][methodName] = stubValue.Method(i)
	}
}

func (s *Server) defaultInvoke(ctx context.Context, req *common.Kmessage) (resp *common.Kmessage, err error) {
	resp.RpcInfo.Timestampe = time.Now().Unix()
	// 解析req 调用反射方法 todo 这个map可能会出问题 其实应该做成map[serviceName][Method]这种形式才比较合适
	if service, ok := s.method[req.RpcInfo.ServiceName]; ok {
		if fn, ok := service[req.RpcInfo.MethodName]; ok {
			ins := []reflect.Value{}
			for i := 0; i < len(req.In.Param); i++ {
				ins = append(ins, reflect.ValueOf(req.In.Param[i]))
			}
			outs := fn.Call(ins)
			changeOutsToInterface := []interface{}{}
			for i := 0; i < len(outs); i++ {
				changeOutsToInterface = append(changeOutsToInterface, outs[i].Interface())
			}
			resp.Out.Response = changeOutsToInterface
			resp.RpcInfo = req.RpcInfo

		} else {
			// 找不到service下面的方法
			errMsg := fmt.Sprintf("Not found the method %s in service %v\n", req.RpcInfo.MethodName, req.RpcInfo.ServiceName)
			resp.RpcInfo = req.RpcInfo
			resp.RpcInfo.Err = errMsg

			return resp, errors.New(errMsg)
		}
	} else {
		// 查不到service
		errMsg := fmt.Sprintf("Not found the service %s in server\n", req.RpcInfo.ServiceName)
		resp.RpcInfo = req.RpcInfo
		resp.RpcInfo.Err = errMsg

		return resp, errors.New(errMsg)
	}
	return resp, nil
}

// 处理解析和反射调用和编码的过程
func (s *Server) serve(conn net.Conn, final common.Endpoint, errChannel chan error) {
	for {
		// 新建req和resp的对象去解析conn里面的req和返回resp
		req := &common.Kmessage{}
		if err := s.remoteServer.DecodeRequest(conn, req); err != nil {
			errChannel <- err
			// 返回给对端的错误信息
			req.RpcInfo.Err = fmt.Sprintf("server decode message error: %v\n", err.Error())
			req.RpcInfo.Timestampe = time.Now().Unix()
			errChannel <- s.remoteServer.EncoceResponse(conn, req)
			continue
		}

		resp, err := s.invoke(context.Background(), req)
		if err != nil {
			errChannel <- err
			// 返回给对端的错误信息
			req.RpcInfo.Err = fmt.Sprintf("server invoke error: %v\n", err.Error())
			req.RpcInfo.Timestampe = time.Now().Unix()
			errChannel <- s.remoteServer.EncoceResponse(conn, req)
			continue
		}

		if err := s.remoteServer.EncoceResponse(conn, resp); err != nil {
			errChannel <- err
			// 返回给对端的错误信息
			req.RpcInfo.Err = fmt.Sprintf("server encode messasge error: %v\n", err.Error())
			req.RpcInfo.Timestampe = time.Now().Unix()
			errChannel <- s.remoteServer.EncoceResponse(conn, req)
			continue
		}
	}

}

// 真正的调用接口，无限循环
func (s *Server) Serve() chan error {
	errChannel := make(chan error, 20)
	// 封装最后的Invoke
	s.invoke = middleware.WrapperMiddleware(s.mds, s.invoke)
	go func(ln net.Listener, c chan error) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				c <- err
				continue
			}
			go s.serve(conn, s.invoke, c)
		}
	}(*s.ls, errChannel)
	return errChannel
}

// 封装一个start接口，确保不用自己Handler错误信息
func (s *Server) Start() {
	errChannel := s.Serve()
	for {
		select {
		case err := <-errChannel:
			klog.Errf("server error in Serve. %v", err.Error())
		default:
			klog.Debugf("listen server error in Serve")
		}
	}
}

// 获取默认配置
func GetDefaultYamlOption() (*option, error) {
	opt := &option{}
	if data, err := ioutil.ReadFile("./default.yaml"); err != nil {
		klog.Infof("can not get yaml in ./default.yaml. next use the default option in code\n")
		return nil, err
	} else {
		var config yamlConfig
		if err := yaml.Unmarshal(data, &config); err != nil {
			klog.Infof("can not unmarshal yaml config in ./default.yaml. next use the default option in code\n")
			return nil, err
		}
		opt.ServerIp = config.Service.Ip
		opt.ServerPort = config.Service.Port
		opt.ServerProtocol = config.Service.Protocol
		opt.psm = config.Service.Psm
		opt.DiscoveryAddress = config.Discovery.Address
		klog.Infof("get the yaml data: %+v\n", opt)
		return opt, nil
	}

}

// 获取指定的Yaml路径配置
func GetYamlOption(path string) (*option, error) {
	opt := &option{}
	if data, err := ioutil.ReadFile(path); err != nil {
		klog.Infof("can not get yaml in %s. next use the default option in code\n", path)
		return nil, err
	} else {
		var config yamlConfig
		if err := yaml.Unmarshal(data, &config); err != nil {
			klog.Infof("can not unmarshal yaml config in %s. next use the default option in code\n", path)
			return nil, err
		}
		opt.ServerIp = config.Service.Ip
		opt.ServerPort = config.Service.Port
		opt.ServerProtocol = config.Service.Protocol
		opt.psm = config.Service.Protocol
		opt.DiscoveryAddress = config.Discovery.Address
		//tod 自定义参数需要添加上
		klog.Infof("get the yaml data: %+v\n", opt)
		return opt, nil
	}
}

// 不读取yaml采用的方式
func CreateOption(port int, ip, protocol, psm string) *option {
	return &option{ServerIp: ip, ServerPort: port, ServerProtocol: protocol, CustomizeProps: map[string]interface{}{}}
}

func NewDefaultServer(opt *option) (*Server, error) {
	server := &Server{}
	defaultListen, err := net.Listen(opt.ServerProtocol, fmt.Sprintf("%s:%d", opt.ServerIp, opt.ServerPort))
	if err != nil {
		klog.Errf("create default listen error: %v", err.Error())
		return nil, err
	}
	server.AddStartHooks(&RegiserHooks{SerivceName: opt.psm, ServiceIp: opt.ServerIp, ServicePort: opt.ServerPort, ConsulAddress: opt.DiscoveryAddress, UseTcp: true})
	server.init(&defaultListen, core.NewDefaultKrpcServer())
	return server, nil
}
