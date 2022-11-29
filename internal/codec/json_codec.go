package codec

import (
	"encoding/json"
	"io"
	"sync"
)

var _ Codec = (*JSONCodec)(nil)

type JSONCodec struct {
	decMap map[io.Reader]*json.Decoder
	encMap map[io.Writer]*json.Encoder
	mutex  *sync.Mutex
}

func NewJSONCodec() *JSONCodec {
	return &JSONCodec{decMap: map[io.Reader]*json.Decoder{}, encMap: map[io.Writer]*json.Encoder{}, mutex: &sync.Mutex{}}
}

// json的deocde形式
func (jsc *JSONCodec) Decode(conn io.Reader, msg interface{}) error {
	var dec *json.Decoder
	if d, ok := jsc.decMap[conn]; ok {
		dec = d
	} else {
		d = json.NewDecoder(conn)
		jsc.mutex.Lock()
		jsc.decMap[conn] = d
		defer jsc.mutex.Unlock()
	}
	if err := dec.Decode(msg); err != nil {
		return err
	}
	return nil
}

func (jsc *JSONCodec) Encode(conn io.Writer, msg interface{}) error {
	var enc *json.Encoder
	if e, ok := jsc.encMap[conn]; ok {
		enc = e
	} else {
		enc = json.NewEncoder(conn)
		jsc.mutex.Lock()
		jsc.encMap[conn] = enc
		defer jsc.mutex.Unlock()
	}
	if err := enc.Encode(msg); err != nil {
		return err
	}
	return nil
}

func (jsc *JSONCodec) Close() error {
	jsc.decMap = map[io.Reader]*json.Decoder{}
	jsc.encMap = map[io.Writer]*json.Encoder{}
	jsc.mutex = nil
	return nil
}
