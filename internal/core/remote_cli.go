package core

import (
	"kRpc/internal/codec"
	"kRpc/internal/common"
	"net"
)

var _ common.RemoteClient = (*KrpcClient)(nil)

type KrpcClient struct {
	co codec.Codec
}

func NewDefaultKrpcClient() *KrpcClient {
	return &KrpcClient{co: codec.NewJSONCodec()}
}

func (r *KrpcClient) Init(co codec.Codec) error {
	r.co = co
	return nil
}

func (r *KrpcClient) EncodeRequest(conn net.Conn, msg *common.Kmessage) error {
	if err := r.co.Encode(conn, msg); err != nil {
		return err
	}
	return nil
}
func (r *KrpcClient) DecodeResponse(conn net.Conn, msg *common.Kmessage) error {
	if err := r.co.Decode(conn, msg); err != nil {
		return err
	}
	return nil

}
func (r *KrpcClient) Close() error {
	return r.co.Close()
}
