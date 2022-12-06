package parse

import (
	"errors"
	"io/ioutil"
	"kRpc/pkg/klog"
	"os"
	"path/filepath"
	"unicode"
)

type ParseError struct {
	desc string
}

func (p *ParseError) Error() string {
	return "parse error: " + p.desc
}

type FieldPair struct {
	Type           string
	SequenceNumber string
}

type Param struct {
	Type string
	Name string
}

type ServiceInfo struct {
	ins  []Param
	outs []Param
}
type KrpcParse struct {
	meta    []byte
	version string
	idx     int
	message map[string]map[string]FieldPair
	service map[string]ServiceInfo
}

func (k *KrpcParse) Parse(path string) error {
	if _, err := os.Stat(path); err != nil {
		klog.Errf("file not exist: %s err: %v\n", path, err.Error())
		return err
	}
	fileSuffix := filepath.Ext(path)
	if fileSuffix != ".krpc" {
		klog.Errf("file suffix is %s instead of .krpc\n", fileSuffix)
		return errors.New("suffix error")
	}
	idl, err := ioutil.ReadFile(path)
	if err != nil {
		klog.Errf("read krpc file %s error: %v\n", path, err.Error())
		return err
	}
	k.meta = idl
	k.parse()
	return nil
}

// 总体的parse函数
func (k *KrpcParse) parse() error {
	// 编译对应版本
	if err := k.parseVersion(); err != nil {
		return err
	}
	// 编译剩余的空格和换行
	k.parseSpace()
	for k.idx < len(k.meta) {
		if k.idx+7 < len(k.meta) {
			if byteSliceEqual(k.meta[k.idx:k.idx+7], "message") {
				k.idx += 7
				if err := k.parseMessage(); err != nil {
					return err
				}
			} else if byteSliceEqual(k.meta[k.idx:k.idx+7], "service") {
				k.idx += 7
				if err := k.parseService(); err != nil {
					return err
				}
			} else {
				return &ParseError{desc: "not found keyword message or service"}
			}
		} else {
			return &ParseError{desc: "not found keyword message or service"}
		}
		k.parseSpace()
	}
	return nil
}

// parse version版本
func (k *KrpcParse) parseVersion() error {
	k.parseSpace()
	if k.idx+7 < len(k.meta) && byteSliceEqual(k.meta[k.idx:k.idx+7], "version") {
		if err := k.parseEqual(); err != nil {
			return err
		}
		if k.idx < len(k.meta) && k.meta[k.idx] == '"' {
			k.idx++
		} else {
			return &ParseError{desc: "version syntax error"}
		}
		start := k.idx
		for k.idx < len(k.meta) && k.meta[k.idx] != '"' {
			k.idx++
		}
		if k.idx >= len(k.meta) {
			return &ParseError{desc: "version syntax error"}
		}
		k.version = string(k.meta[start:k.idx])
		k.idx++
		return nil
	}
	return &ParseError{desc: "miss version description"}
}

func (k *KrpcParse) parseEqual() error {
	k.parseSpace()
	if k.idx < len(k.meta) && k.meta[k.idx] == '=' {
		return nil
	}
	return &ParseError{desc: "miss ="}
}

// 解析序号
func (k *KrpcParse) parseSequenceNumber() (string, error) {
	k.parseSpace()
	start := k.idx
	for k.idx < len(k.meta) && unicode.IsDigit(rune(k.meta[k.idx])) {
		k.idx++
	}
	if start == k.idx {
		return "", &ParseError{desc: "miss seuqenceNumber in field"}
	}
	return string(k.meta[start:k.idx]), nil
}

// 解析类型
func (k *KrpcParse) parseType() (string, error) {
	return "", nil
}

// 解析名字
func (k *KrpcParse) parseName() (string, error) {
	k.parseSpace()
	start := k.idx
	for k.idx < len(k.meta) && k.meta[k.idx] != ' ' && k.meta[k.idx] != '=' {
		k.idx++
	}
	if k.idx >= len(k.meta) {
		return "", &ParseError{desc: "parse token error, the variable name parse error"}
	}
	return string(k.meta[start:k.idx]), nil

}

// parse message 结构 {string tag1 = 1; string tag2 = 2;} todo 后面需要加上一个状态机
func (k *KrpcParse) parseMessage() error {
	k.parseSpace()
	if k.meta[k.idx] != '{' {
		return &ParseError{desc: "{ should be append after message"}
	}
	k.idx++
	for k.idx < len(k.meta) && k.meta[k.idx] != '}' {
		k.parseType()
		k.parseName()
		k.parseEqual()
		k.parseSequenceNumber()
		k.parseDot()
	}
	if k.meta[k.idx] != '}' || k.idx >= len(k.meta) {
		return &ParseError{desc: "} should warp the message"}
	}
	k.idx++
	k.parseSpace()
	return nil
}

// parse service
func (k *KrpcParse) parseService() error {
	return nil
}

// parse \n and ' '
func (k *KrpcParse) parseSpace() {
	for k.idx < len(k.meta) && (k.meta[k.idx] == ' ' || k.meta[k.idx] == '\n' || k.meta[k.idx] == 'r') {
		k.idx++
	}
}

// parse ;
func (k *KrpcParse) parseDot() {
	k.parseSpace()
	for k.idx < len(k.meta) && k.meta[k.idx] == ';' {
		k.idx++
		k.parseSpace()
	}
}

func byteSliceEqual(a []byte, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
