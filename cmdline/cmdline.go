package main

import (
	"cmdline/parse"
	"klog"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	kp := parse.NewKrpcParse()
	var path string
	var cmdGenerater = &cobra.Command{
		Use:   "parse",
		Short: "parse and generate",
		Long:  "parse xxx.krpc to generate go's code to use",
		Run: func(cmd *cobra.Command, args []string) {
			kp.Parse(path)
		},
	}
	cmdGenerater.Flags().StringVarP(&path, "path", "p", "", "krpc file's path")

	if err := cmdGenerater.Execute(); err != nil {
		klog.Errf("cmdline error: %v\n", err.Error())
		os.Exit(1)
	}
}
