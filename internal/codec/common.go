package codec

import "io"

// 编码通用实现
type Codec interface {
	Decode(conn io.Reader, msg interface{}) error
	Encode(conn io.Writer, msg interface{}) error
	Close() error
}
