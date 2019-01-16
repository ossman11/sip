// +build js

package main

import (
	"os"
	"syscall/js"
	"testing"

	"github.com/ossman11/sip/core/def"
)

func TestMain(m *testing.M) {
	def.Integration()
	result := m.Run()
	os.Exit(result)
}

func TestCore(t *testing.T) {
	main()

	global := js.Global()
	sipJS := global.Get("sip")

	if sipJS.Type() == js.TypeUndefined {
		t.Error("Failed to lookup sip Object as type is undefined.")
	}
}
