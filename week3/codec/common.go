package codec

import "io"

/**
* 通用codec实现，规定为一个decode和encode,以及需要保留的option
**/

type Codec interface {
	Decode(conn io.Reader, resp interface{}) error
	Encode(conn io.Writer, req interface{}) error
}
