package main

import (
	"context"
	"fmt"
)

type endpoint func(ctx context.Context, req interface{}) (resp interface{}, err error)

type middleware func(ctx context.Context, req interface{}, handler endpoint) (resp interface{}, err error)

type warp func(endpoint) endpoint

func main() {
	mds := []middleware{}
	mds = append(mds, func(ctx context.Context, req interface{}, handler endpoint) (resp interface{}, err error) {
		fmt.Printf("before1\n")
		resp, err = handler(ctx, req)
		fmt.Printf("end1\n")
		return
	})
	mds = append(mds, func(ctx context.Context, req interface{}, handler endpoint) (resp interface{}, err error) {
		fmt.Printf("before2\n")
		resp, err = handler(ctx, req)
		fmt.Printf("end2\n")
		return
	})
	var handler endpoint = func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		fmt.Printf("make msg\n")
		return nil, nil
	}
	for i := len(mds) - 1; i >= 0; i-- {
		handler = warp(func(e endpoint) endpoint {
			// 由于go的机制问题如果不用tmp去存下当前的i，那么mds[i]就会取最终的那一个，就会溢出，所以在return前先保存一下i的量，然后每一个stack去存的变量就是对的
			cur := i
			return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
				return mds[cur](ctx, req, e)
			}
		})(handler)
	}
	resp, err := handler(context.Background(), "ster")
	if resp != nil && err != nil {
		return
	}
}
