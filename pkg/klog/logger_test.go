package klog

import (
	"testing"
)

func TestInfof(t *testing.T) {
	t.Run("info", func(t *testing.T) {
		Infof("%s\n", "zhaohaokai")
	})
}
