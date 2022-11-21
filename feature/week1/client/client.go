package client

import (
	"fmt"
	"net"
)

type krpcClient struct {
}

func (k *krpcClient) Call(msg interface{}) {
	conn, err := net.Dial("tcp", ":10011")
	buffer := make([]byte, 1024)
	if err != nil {
		panic("[kRpc Client] connect to the server error!")
	}
	if _, err := conn.Write([]byte("I am haokai")); err != nil {
		fmt.Printf("[kRpc Client] call remote program error: %v\n", err.Error())
	}
	idx, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("[kRpc Client] read remote program response: %v", err.Error())
		fmt.Println()
	}
	fmt.Printf("[kRpc Client] the result from server is %s ", string(buffer[:idx]))
	fmt.Println()
}

func NewKrpcClient() *krpcClient {
	return &krpcClient{}
}
