package middleware

import (
	"context"
	"kRpc/internal/common"
)

type Middleware func(ctx context.Context, req *common.Kmessage, handler common.Endpoint) (resp *common.Kmessage, err error)

type mdsWraper func(in common.Endpoint) (out common.Endpoint)

func WrapperMiddleware(mds []Middleware, handler common.Endpoint) common.Endpoint {
	for i := len(mds) - 1; i >= 0; i-- {
		handler = mdsWraper(func(in common.Endpoint) common.Endpoint {
			cur := i
			return func(ctx context.Context, req *common.Kmessage) (resp *common.Kmessage, err error) {
				return mds[cur](ctx, req, in)
			}
		})(handler)
	}
	return handler
}
