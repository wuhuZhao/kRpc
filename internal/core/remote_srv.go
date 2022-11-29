package core

import (
	"kRpc/internal/codec"
	"kRpc/internal/common"
	"net"
)

var _ common.RemoteServer = (*KrpcServer)(nil)

type KrpcServer struct {
	co codec.Codec
}

func NewDefaultKrpcServer() *KrpcServer {
	return &KrpcServer{co: codec.NewJSONCodec()}
}

func (r *KrpcServer) Init(co codec.Codec) error {
	r.co = co
	return nil
}

func (r *KrpcServer) DecodeRequest(conn net.Conn, msg *common.Kmessage) error {
	if err := r.co.Decode(conn, msg); err != nil {
		return err
	}
	return nil

}

func (r *KrpcServer) EncoceResponse(conn net.Conn, msg *common.Kmessage) error {
	if err := r.co.Encode(conn, msg); err != nil {
		return err
	}
	return nil
}

func (r *KrpcServer) Close() error {
	return r.co.Close()
}
