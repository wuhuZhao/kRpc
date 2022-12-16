package main

import "context"

type MyService interface {
	getCid(ctx context.Context, cid int32, cid2 Req) (cid1 *Req, err error)
}
type MyServiceImpl struct{}

var _ MyService = (*MyServiceImpl)(nil)

type Req struct {
	cid int32

	cid2 int64
}

func (impl *MyServiceImpl) getCid(ctx context.Context, cid int32, cid2 Req) (cid1 *Req, err error) {
	return nil, nil
}
