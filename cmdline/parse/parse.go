package parse

import (
	"errors"
	"fmt"
	"io/ioutil"
	"klog"
	"os"
	"path/filepath"
	"unicode"
)

var TypeEnum = []string{"string", "int32", "int64", "float", "void"}

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

func (k *KrpcParse) ToPrint() {
	klog.Infof("meta: %s \nversion: %s \nmessage: %v \nservice: %v \nidx: %d", string(k.meta), k.version, k.message, k.service, k.idx)
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
	return k.parse()
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
				messageName, err := k.parseBlockName()
				if err != nil {
					return err
				}
				if err := k.parseMessage(messageName); err != nil {
					return err
				}
			} else if byteSliceEqual(k.meta[k.idx:k.idx+7], "service") {
				k.idx += 7
				serviceName, err := k.parseBlockName()
				if err != nil {
					return err
				}
				if err := k.parseService(serviceName); err != nil {
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

// parse message 后面的名字
func (k *KrpcParse) parseBlockName() (string, error) {
	k.parseSpace()
	start := k.idx
	for k.idx < len(k.meta) {
		if k.meta[k.idx] != ' ' && k.meta[k.idx] != '}' {
			k.idx++
		} else {
			break
		}
	}
	return string(k.meta[start:k.idx]), nil
}

// parse version版本
func (k *KrpcParse) parseVersion() error {
	k.parseSpace()
	if k.idx+7 < len(k.meta) && byteSliceEqual(k.meta[k.idx:k.idx+7], "version") {
		k.idx += 7
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
		k.idx++
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
	k.parseSpace()
	start := k.idx
	for k.idx < len(k.meta) && k.meta[k.idx] != ' ' {
		k.idx++
	}
	for i := 0; i < len(TypeEnum); i++ {
		if byteSliceEqual(k.meta[start:k.idx], TypeEnum[i]) {
			return string(k.meta[start:k.idx]), nil
		}
	}
	for messageName := range k.message {
		if byteSliceEqual(k.meta[start:k.idx], messageName) {
			return string(k.meta[start:k.idx]), nil
		}
	}
	return "", &ParseError{desc: fmt.Sprintf("%s not found the type in %v and %v", k.meta[start:k.idx], TypeEnum, k.message)}
}

// 解析名字
func (k *KrpcParse) parseName() (string, error) {
	k.parseSpace()
	start := k.idx
	for k.idx < len(k.meta) && k.meta[k.idx] != ' ' && k.meta[k.idx] != '=' && k.meta[k.idx] != ')' {
		k.idx++
	}
	if k.idx >= len(k.meta) {
		return "", &ParseError{desc: "parse token error, the variable name parse error"}
	}
	return string(k.meta[start:k.idx]), nil

}

// parse message 结构 {string tag1 = 1; string tag2 = 2;} todo 后面需要加上一个状态机
func (k *KrpcParse) parseMessage(messageName string) error {
	k.parseSpace()
	if k.meta[k.idx] != '{' {
		return &ParseError{desc: "{ should be append after message"}
	}
	k.idx++
	for k.idx < len(k.meta) && k.meta[k.idx] != '}' {
		t, err := k.parseType()
		if err != nil {
			return err
		}
		n, err := k.parseName()
		if err != nil {
			return err
		}
		if err := k.parseEqual(); err != nil {
			return err
		}
		s, err := k.parseSequenceNumber()
		if err != nil {
			return err
		}
		if err := k.parseDot(); err != nil {
			return err
		}
		if k.message[messageName] == nil {
			k.message[messageName] = map[string]FieldPair{}
		}
		k.message[messageName][n] = FieldPair{Type: t, SequenceNumber: s}
	}
	if k.meta[k.idx] != '}' || k.idx >= len(k.meta) {
		return &ParseError{desc: "} should warp the message"}
	}
	k.idx++
	k.parseSpace()
	return nil
}

// parse service service {rpc getResp(Req req) return Resp}
func (k *KrpcParse) parseService(serviceName string) error {
	k.parseSpace()
	if k.meta[k.idx] != '{' {
		return &ParseError{desc: "{ should be append after message"}
	}
	k.idx++
	for k.idx < len(k.meta) && k.meta[k.idx] != '}' {
		funcName, err := k.parseFuncName()
		if err != nil {
			return err
		}
		ins, err := k.parseParams()
		if err != nil {
			return err
		}
		if err := k.skipReturn(); err != nil {
			return err
		}
		outs, err := k.parseParams()
		if err != nil {
			return err
		}
		k.service[funcName] = ServiceInfo{ins: ins, outs: outs}
		k.parseSpace()
	}
	if k.idx >= len(k.meta) || k.meta[k.idx] != '}' {
		return &ParseError{desc: "} should warp the message"}
	}
	k.idx++
	k.parseSpace()
	return nil
}

func (k *KrpcParse) skipReturn() error {
	k.parseSpace()
	if k.idx+6 >= len(k.meta) || !byteSliceEqual(k.meta[k.idx:k.idx+6], "return") {
		return &ParseError{desc: "miss return in service"}
	}
	k.idx += 6
	return nil
}

func (k *KrpcParse) parseParams() ([]Param, error) {
	k.parseSpace()
	if k.idx+1 >= len(k.meta) || k.meta[k.idx] != '(' {
		return nil, &ParseError{desc: "( shoule be append after function name or returns field"}
	}
	k.idx++
	res := []Param{}
	for k.idx < len(k.meta) {
		k.parseSpace()
		t, err := k.parseType()
		if err != nil {
			return nil, err
		}
		n, err := k.parseName()
		if err != nil {
			return nil, err
		}
		res = append(res, Param{Type: t, Name: n})
		k.parseSpace()
		klog.Errf("test %s\n", string(k.meta[k.idx]))
		if k.idx < len(k.meta) && k.meta[k.idx] == ')' {
			break
		} else if k.idx < len(k.meta) && k.meta[k.idx] == ',' {
			k.idx++
			continue
		}
	}
	if k.idx >= len(k.meta) || k.meta[k.idx] != ')' {
		return nil, &ParseError{desc: "( shoule be append after function name or returns field"}
	}
	k.idx++
	k.parseSpace()
	return res, nil
}

func (k *KrpcParse) parseFuncName() (string, error) {
	k.parseSpace()
	if k.idx+3 >= len(k.meta) || !byteSliceEqual(k.meta[k.idx:k.idx+3], "rpc") {
		return "", &ParseError{desc: "syntax error not found 'rpc'"}
	}
	k.idx += 3
	k.parseSpace()
	start := k.idx
	for k.idx < len(k.meta) {
		if k.meta[k.idx] != '(' && k.meta[k.idx] != ' ' {
			k.idx++
		} else {
			break
		}
	}
	return string(k.meta[start:k.idx]), nil
}

// parse \n and ' '
func (k *KrpcParse) parseSpace() {
	for k.idx < len(k.meta) && (k.meta[k.idx] == ' ' || k.meta[k.idx] == '\n' || k.meta[k.idx] == '\r') {
		k.idx++
	}
}

// parse ;
func (k *KrpcParse) parseDot() error {
	k.parseSpace()
	flag := false
	for k.idx < len(k.meta) && k.meta[k.idx] == ';' {
		k.idx++
		flag = true
		k.parseSpace()
	}
	if !flag {
		return &ParseError{desc: "miss ';' in the message field"}
	}
	return nil
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

func NewKrpcParse() *KrpcParse {
	return &KrpcParse{
		meta:    []byte{},
		idx:     0,
		message: map[string]map[string]FieldPair{},
		service: map[string]ServiceInfo{},
	}
}
