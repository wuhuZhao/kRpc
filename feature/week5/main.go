package main

import (
	"context"
	"fmt"
	"time"
	"week5/client"
	"week5/core"
	"week5/server"
)

type HelloWorldServiceImpl struct{}

func (h *HelloWorldServiceImpl) GetPsmInfo(cluster string) string {
	return "data.byte.diff" + ":" + cluster
}

// call执行测试test
func call() {
	time.Sleep(5 * time.Second)
	cli, err := client.NewDefaultClient(&client.Option{ServerProtocol: "tcp", ServiceName: "krpcServer"})
	if err != nil {
		fmt.Printf("client create error: %v\n", err)
	}
	// 插入中间件  打印参数和结果
	cli.Use(func(e core.Endpoint) core.Endpoint {
		return func(ctx context.Context, req, resp *core.Message) error {
			fmt.Printf("client req : %v \n", req.RpcInfo)
			return e(ctx, req, resp)
		}
	})
	request := &core.Message{RpcInfo: &core.RpcInfo{ServiceName: "HelloWorldServiceImpl", MethodName: "GetPsmInfo", Param: []interface{}{"haokaizhao"}}}
	if _, err := cli.Call(request); err != nil {
		fmt.Printf("client call error: %v\n", err)
	}
	request = &core.Message{RpcInfo: &core.RpcInfo{ServiceName: "HelloWorldServiceImpl", MethodName: "GetPsmInfo", Param: []interface{}{"jasonbsun"}}}
	if _, err := cli.Call(request); err != nil {
		fmt.Printf("client call error: %v\n", err)
	}
}

func main() {
	srv, err := server.NewDefaultServer("./example.yaml")
	if err != nil {
		fmt.Printf("server crete error: %v\n", err)
	}
	// 插入中间件打印request和response
	srv.Use(func(e core.Endpoint) core.Endpoint {
		return func(ctx context.Context, req, resp *core.Message) error {
			fmt.Printf("server req: %v\n", req.RpcInfo)
			return e(ctx, req, resp)
		}
	})
	srv.Register(&HelloWorldServiceImpl{})
	errchan := srv.Serve()
	go call()
hang:
	select {
	case _ = <-errchan:
		goto hang
	}
}
