// package main

// import (
// 	"context"
// 	"fmt"
// )

// type ep func(context.Context, struct{}, struct{}) error

// type next func(ep) ep

// func (t *test) end(ctx context.Context, req struct{}, resp struct{}) error {
// 	fmt.Println("end")
// 	return nil
// }

// // n里面实现自己的逻辑
// func (t *test) Use(n next) {
// 	t.final = (t.final)
// }

// type test struct {
// 	final ep
// 	next  next
// }

// func (t *test) init() {
// 	t.final = t.end
// }
