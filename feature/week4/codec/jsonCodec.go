package codec

import (
	"encoding/json"
	"io"
)

// 保证JsonCodec实现Codec
var _ Codec = (*JsonCodec)(nil)

// json的option配置
type Option struct {
	filter       map[string]struct{}
	defaultValue map[string]interface{}
}

// json的codec具体实现
type JsonCodec struct {
	option *Option
}

func (jsc *JsonCodec) Decode(conn io.Reader, resp interface{}) error {
	dec := json.NewDecoder(conn)
	err := dec.Decode(resp)
	if err != nil {
		return err
	}
	return nil
}

func (jsc *JsonCodec) Encode(conn io.Writer, req interface{}) error {
	enc := json.NewEncoder(conn)
	err := enc.Encode(req)
	if err != nil {
		return err
	}
	return nil
}

func NewJsonCodec(opt *Option) *JsonCodec {
	return &JsonCodec{option: opt}
}
