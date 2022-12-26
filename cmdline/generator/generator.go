package generator

import (
	"bytes"
	"cmdline/parse"
	"io/fs"
	"io/ioutil"
	"strings"
	"text/template"
)

type krpcGenerator struct {
	message           map[string]map[string]parse.FieldPair
	service           map[string]parse.ServiceInfo
	messageTemplate   string
	serviceTemplate   string
	interfaceTemplate string
	serviceName       string
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
func (k *krpcGenerator) Generate() error {
	tpl := template.New("tmpl")
	code := &strings.Builder{}
	if err := k.generateInterface(code, tpl); err != nil {
		return err
	}
	code.WriteByte('\n')
	if err := k.generateMessage(code, tpl); err != nil {
		return err
	}
	code.WriteByte('\n')
	if err := k.generatService(code, tpl); err != nil {
		return err
	}
	if err := ioutil.WriteFile("./generator_krpc.go", []byte(code.String()), fs.ModeAppend); err != nil {
		return err
	}
	return nil
}

func (k *krpcGenerator) generatService(code *strings.Builder, tpl *template.Template) error {
	serviceTpl, err := tpl.Parse(krpcService())
	if err != nil {
		return err
	}
	info := map[string]string{
		"ServiceName": k.serviceName,
	}
	for funcName, serviceInfo := range k.service {
		info["FuncName"] = funcName
		info["FuncParam"] = k.generateParams(serviceInfo.Ins)
		info["FuncType"] = k.generateParams(serviceInfo.Outs)
		info["ReturnNil"] = func(returns []parse.Param) string {
			data := &strings.Builder{}
			for i := 0; i < len(returns); i++ {
				data.WriteString(fillNilValue(returns[i].Type))
				if i != len(returns)-1 {
					data.WriteByte(',')
				}
			}
			return data.String()
		}(serviceInfo.Outs)
		var res bytes.Buffer
		if err := serviceTpl.Execute(&res, info); err != nil {
			return err
		}
		code.WriteString(res.String())
		code.WriteByte('\n')
	}
	return nil
}

func (k *krpcGenerator) generateMessage(code *strings.Builder, tpl *template.Template) error {
	messageTpl, err := tpl.Parse(krpcMessage())
	if err != nil {
		return err
	}
	for messageName, messageInfo := range k.message {
		info := map[string]interface{}{}
		info["StructName"] = messageName
		for fieldName, FieldInfo := range messageInfo {
			info["StructFieldList"] = append(info["StructFieldList"].([]map[string]string),
				map[string]string{
					"Field": fieldName,
					"Type":  FieldInfo.Type,
				})
		}
		var res bytes.Buffer
		if err := messageTpl.Execute(&res, info); err != nil {
			return err
		}
		code.WriteString(res.String())
		code.WriteByte('\n')
	}
	return nil
}

func (k *krpcGenerator) generateInterface(code *strings.Builder, tpl *template.Template) error {
	interfaceTpl, err := tpl.Parse(krpcInterface())
	if err != nil {
		return err
	}
	info := map[string]interface{}{
		"ServiceName": k.serviceName,
	}
	for funcName, serviceInfo := range k.service {
		info["InterfaceList"] = append(info["InterfaceList"].([]map[string]string), map[string]string{
			"FuncName":  funcName,
			"FuncParam": k.generateParams(serviceInfo.Ins),
			"FuncType":  k.generateParams(serviceInfo.Outs),
		})
	}
	var res bytes.Buffer
	if err := interfaceTpl.Execute(&res, info); err != nil {
		return err
	}
	code.WriteString(res.String())
	return nil
}

func (k *krpcGenerator) generateParams(p []parse.Param) string {
	res := &strings.Builder{}

	if len(p) == 1 {
		return p[0].Name + " " + p[0].Type
	}
	for i := 0; i < len(p); i++ {
		res.WriteString(p[i].Name)
		res.WriteByte(' ')
		res.WriteString(p[i].Type)
		if i == len(p)-1 {
			continue
		}
		res.WriteByte(',')
	}
	return res.String()
}

func fillNilValue(key string) string {
	switch key {
	case "int32", "int64":
		return "0"
	case "string":
		return "\"\""
	case "float":
		return "0.0"
	default:
		return "nil"
	}
}
