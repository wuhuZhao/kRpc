package client

import (
	"fmt"
	"net"
	"week3/codec"
	"week3/core"
)

type krpcClient struct {
	co codec.Codec
}

func (k *krpcClient) Call(msg *core.Message) {
	conn, err := net.Dial("tcp", ":10011")
	if err != nil {
		fmt.Printf("[kRpcClient] eror connect %v\n", err)
	}
	if er := k.co.Encode(conn, msg); er != nil {
		fmt.Printf("[kRpcClient] error encode %v\n", er)
	}
	var resp interface{}
	if er := k.co.Decode(conn, &resp); er != nil {
		fmt.Printf("[kRpcClient] error decode %v\n", er)
	}
	fmt.Printf("[kRpcClient] resp %v\n", resp)
}

func NewKrpcClient() *krpcClient {
	return &krpcClient{co: codec.NewJsonCodec(nil)}
}
