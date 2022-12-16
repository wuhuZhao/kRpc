package generator

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"testing"
	"text/template"
)

func TestGenerateGo(t *testing.T) {
	t2 := template.New("tmpl")
	goCode := `package main
	 import "context"`

	if t3, err := t2.Parse(krpcInterface()); err != nil {
		t.Fatal(err.Error())
	} else {
		var final bytes.Buffer
		m := map[string]interface{}{
			"ServiceName": "MyService",
			"InterfaceList": []map[string]string{
				map[string]string{
					"FuncName":  "getCid",
					"FuncParam": "ctx context.Context, cid int32, cid2 Req",
					"FuncType":  "cid1 *Req, err error",
				},
			},
		}
		if err := t3.Execute(&final, m); err != nil {
			t.Fatal(err.Error())
		}
		goCode += final.String()
		t.Log(goCode)
	}
	if t3, err := t2.Parse(krpcMessage()); err != nil {
		t.Fatal(err.Error())
	} else {
		var final bytes.Buffer
		m := map[string]interface{}{}
		m["StructName"] = "Req"
		m["StructFieldList"] = []map[string]string{}
		m["StructFieldList"] = append(m["StructFieldList"].([]map[string]string), map[string]string{
			"Field": "cid",
			"Type":  "int32",
		})
		m["StructFieldList"] = append(m["StructFieldList"].([]map[string]string), map[string]string{
			"Field": "cid2",
			"Type":  "int64",
		})
		if err := t3.Execute(&final, m); err != nil {
			t.Fatal(err.Error())
		}
		goCode += "\n" + final.String()
		t.Log(goCode)
	}
	if t3, err := t2.Parse(krpcService()); err != nil {
		t.Fatal(err.Error())
	} else {
		var final bytes.Buffer
		m := map[string]string{
			"ServiceName": "MyService",
			"FuncName":    "getCid",
			"FuncParam":   "ctx context.Context, cid int32, cid2 Req",
			"FuncType":    "cid1 *Req, err error",
			"ReturnNil":   "nil, nil",
		}
		if err := t3.Execute(&final, m); err != nil {
			t.Fatal(err.Error())
		}
		goCode += "\n" + final.String()
		t.Log(goCode)
	}
	ioutil.WriteFile("../generator_krpc.go", []byte(goCode), fs.ModeAppend)
}

/**
	return `
	type {{.ServiceName}} interface {
		{{- range .InterfaceList}}
		{{.FuncName}} ({{.FuncParam}}) {{.FuncType}}
		{{end}}
	}
	type {{.ServiceName}}Impl struct {}
	var _ {{.ServiceName}} = (*{{.ServiceName}}Impl)(nil)
	`
**/
func TestGenerateInterface(t *testing.T) {
	t2, err := template.New("tmpl").Parse(krpcInterface())
	if err != nil {
		t.Fatal(err.Error())
	}
	var final bytes.Buffer
	m := map[string]string{
		"ServiceName": "MyService",
		"FuncName":    "getCid",
		"FuncParam":   "ctx Context.context, cid int32, cid2 int64",
		"FuncType":    "cid1 int32, err error",
		"ReturnNil":   "nil, nil",
	}
	if err := t2.Execute(&final, m); err != nil {
		t.Fatal(err.Error())
	}
	t.Log(final.String())

}

/**
	return `func (impl *{{.ServiceName}}Impl) {{.FuncName}} ({{.FuncParam}}) ({{.FuncType}})  {
		return {{.ReturnNil}}
	}`
**/
func TestGenerateService(t *testing.T) {
	t2, err := template.New("tmpl").Parse(krpcService())
	if err != nil {
		t.Fatal(err.Error())
	}
	var final bytes.Buffer
	m := map[string]string{
		"ServiceName": "MyService",
		"FuncName":    "getCid",
		"FuncParam":   "ctx Context.context, cid int32, cid2 int64",
		"FuncType":    "cid1 int32, err error",
		"ReturnNil":   "nil, nil",
	}
	if err := t2.Execute(&final, m); err != nil {
		t.Fatal(err.Error())
	}
	t.Log(final.String())
}

/**
	return `type {{.StructName}} struct {
		{{- range .StructFieldList}}
		{{.Field}} {{.Type}}
		{{end}}
	}`
**/
func TestGenerateMessage(t *testing.T) {
	t2, err := template.New("tmpl").Parse(krpcMessage())
	if err != nil {
		t.Fatal(err.Error())
	}
	var final bytes.Buffer
	m := map[string]interface{}{}
	m["StructName"] = "Req"
	m["StructFieldList"] = []map[string]string{}
	m["StructFieldList"] = append(m["StructFieldList"].([]map[string]string), map[string]string{
		"Field": "cid",
		"Type":  "int32",
	})
	m["StructFieldList"] = append(m["StructFieldList"].([]map[string]string), map[string]string{
		"Field": "cid2",
		"Type":  "int64",
	})
	m["StructFieldList"] = append(m["StructFieldList"].([]map[string]string), map[string]string{
		"Field": "cid1",
		"Type":  "Response",
	})
	if err := t2.Execute(&final, m); err != nil {
		t.Fatal(err.Error())
	}
	t.Log(final.String())
}
