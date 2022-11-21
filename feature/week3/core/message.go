package core

type Status int32

const (
	success Status = iota
	fail
	unknow
)

type Message struct {
	*RpcInfo
	Status Status
}

type RpcInfo struct {
	// 路由名
	ServiceName string
	// interface名
	MethodName string
	// 协议名
	Protocol string
	// 传入的参数 server端可以不返回，client端如果调用有参函数则需要返回
	Param []interface{}
	// 传出的结果
	Response []interface{}
	// 此次rpc的时间
	Timestampe int64
	// 此次rpc的上下文id
	Cid int64
	// error
	Err string
}

func NewWithSuccessMessage(r *RpcInfo) *Message {
	return &Message{r, success}
}

func NewWithFailMessage(r *RpcInfo) *Message {
	return &Message{r, fail}
}

func NewWithUnknowMessage(r *RpcInfo) *Message {
	return &Message{r, unknow}
}
