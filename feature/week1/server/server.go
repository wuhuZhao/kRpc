package server

import (
	"fmt"
	"net"
)

type krpcServer struct {
}

func (k *krpcServer) Serve(c chan string) {
	listen, err := net.Listen("tcp", ":10011")
	if err != nil {
		panic("[kRpc] the server can't start.Because the port  is used")
	}
	c <- "ok"
	read := make([]byte, 1024)
	for {
		conn, err := listen.Accept()
		if err != nil {
			panic("[kRpc] the server is work on a error mode when listenning")
		}
		if idx, err := conn.Read(read); err != nil {
			fmt.Printf("[kRpc] the server is work on a error mode when reading. %v s", err)
			fmt.Println()
			continue
		} else {
			conn.Write([]byte(fmt.Sprintf("[kRpc] server is accepted! %s\n", string(read[:idx]))))
		}
		fmt.Println("[kRpc] the server has complete a call")
	}
}

func NewKrpcServer() *krpcServer {
	return &krpcServer{}
}
