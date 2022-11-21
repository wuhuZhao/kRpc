package main

import (
	"fmt"
	"time"
	"week3/client"
	"week3/core"
	"week3/server"
)

type HelloWorldServiceImpl struct{}

func (h *HelloWorldServiceImpl) GetPsmInfo(cluster string) string {
	return "data.byte.diff" + ":" + cluster
}

func call() {
	time.Sleep(5 * time.Second)
	cli := client.NewKrpcClient()
	cli.Call(&core.Message{RpcInfo: &core.RpcInfo{ServiceName: "HelloWorldServiceImpl", MethodName: "GetPsmInfo", Param: []interface{}{"haokaizhao"}}})
}

func main() {
	srv := server.NewKrpcServer(server.NewKrpcOptionWithDefaultMode("127.0.0.1", "10011", "tcp"))
	srv.Register(&HelloWorldServiceImpl{})
	err := srv.Serve()
	go call()
	select {
	case er := <-err:
		{
			fmt.Printf("[krpcServer] connect error %+v", er)
		}
	}

}
