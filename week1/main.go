package main

import (
	"week1/client"
	"week1/server"
)

func main() {
	srv := server.NewKrpcServer()
	inform := make(chan string)
	// 暂时先用channel做通知通信，后面再考虑拆分成conn_init和serve两步
	go srv.Serve(inform)
	// 等待服务端启动
	<-inform
	cli := client.NewKrpcClient()
	cli.Call("hello word!")
}
