package core

import (
	"net"
	"week5/codec"
)

var _ RemoteClient = (*krpcClient)(nil)

type krpcClient struct {
	co codec.Codec
}

func (k *krpcClient) Call(conn net.Conn, request *Message, response *Message) error {
	if er := k.co.Encode(conn, request); er != nil {
		return er
	}
	if er := k.co.Decode(conn, response); er != nil {
		return er
	}
	return nil
}

func NewKrpcClient() *krpcClient {
	return &krpcClient{co: codec.NewJsonCodec(nil)}
}
