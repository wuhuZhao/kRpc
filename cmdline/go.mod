module cmdline

go 1.18

require (
	github.com/spf13/cobra v1.6.1
	klog v1.0.0
)

require (
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace klog v1.0.0 => ../pkg/klog
