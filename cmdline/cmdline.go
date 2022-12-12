package main

import (
	"cmdline/parse"
	"klog"
)

// syntax的parse 和 generator
func main() {
	kp := parse.NewKrpcParse()
	if err := kp.Parse("./test.krpc"); err != nil {
		klog.Errf("err: %s", err.Error())
	}
	kp.ToPrint()
}
