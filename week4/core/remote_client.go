package core

import (
	"fmt"
	"net"
	"week4/codec"
)

var _ RemoteClient = (*krpcClient)(nil)

type krpcClient struct {
	co codec.Codec
}

func (k *krpcClient) Call(conn net.Conn, ep Endpoint, request *Message) (*Message, error) {
	if er := k.co.Encode(conn, request); er != nil {
		fmt.Printf("[kRpcClient] error encode %v\n", er)
		return nil, er
	}
	var resp *Message
	if er := k.co.Decode(conn, resp); er != nil {
		fmt.Printf("[kRpcClient] error decode %v\n", er)
		return nil, er
	}
	return resp, nil
}

func NewKrpcClient() *krpcClient {
	return &krpcClient{co: codec.NewJsonCodec(nil)}
}
