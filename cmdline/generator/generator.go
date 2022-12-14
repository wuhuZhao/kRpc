package generator

import (
	"cmdline/parse"
	"io/ioutil"
	"klog"
)

type krpcGenerator struct {
	message         map[string]map[string]parse.FieldPair
	service         map[string]parse.ServiceInfo
	messageTemplate string
	serviceTemplate string
}

func NewKrpcGnerator(msgPath, servicePath string, message map[string]map[string]parse.FieldPair, service map[string]parse.ServiceInfo) *krpcGenerator {
	msgT, err := ioutil.ReadFile(msgPath)
	if err != nil {
		klog.Errf("msgTemplate path : %s error: %v\n", msgPath, err.Error())
		panic("msgTemplate path error")
	}
	serviceT, err := ioutil.ReadFile(servicePath)
	if err != nil {
		klog.Errf("serviceT path : %s error: %v\n", servicePath, err.Error())
		panic("serviceTemplate path error")
	}
	return &krpcGenerator{
		messageTemplate: string(msgT),
		serviceTemplate: string(serviceT),
		message:         message,
		service:         service,
	}
}

// 循环解析token生成模板
func (k *krpcGenerator) Generate() {

}
