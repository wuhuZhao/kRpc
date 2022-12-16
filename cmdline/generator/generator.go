package generator

import (
	"cmdline/parse"
)

type krpcGenerator struct {
	message           map[string]map[string]parse.FieldPair
	service           map[string]parse.ServiceInfo
	messageTemplate   string
	serviceTemplate   string
	interfaceTemplate string
}

func NewKrpcGnerator(msgPath, servicePath string, message map[string]map[string]parse.FieldPair, service map[string]parse.ServiceInfo) *krpcGenerator {

	return &krpcGenerator{
		messageTemplate:   krpcMessage(),
		serviceTemplate:   krpcMessage(),
		interfaceTemplate: krpcInterface(),
		message:           message,
		service:           service,
	}
}

func krpcInterface() string {
	return `
	type {{.ServiceName}} interface {
		{{- range .InterfaceList}}
		{{.FuncName}} ({{.FuncParam}}) ({{.FuncType}})
		{{end}}
	}
	type {{.ServiceName}}Impl struct {}
	var _ {{.ServiceName}} = (*{{.ServiceName}}Impl)(nil)
	`
}

func krpcService() string {
	return `func (impl *{{.ServiceName}}Impl) {{.FuncName}} ({{.FuncParam}}) ({{.FuncType}})  {
		return {{.ReturnNil}}
	}`
}

func krpcMessage() string {
	return `type {{.StructName}} struct {
		{{- range .StructFieldList}}
		{{.Field}} {{.Type}}
		{{end}}
	}`
}

// 循环解析token生成模板
func (k *krpcGenerator) Generate() {

}
