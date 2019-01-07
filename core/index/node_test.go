package index

import (
	"net"
	"os"
	"testing"

	"github.com/ossman11/sip/core/def"
)

func Test_resolveLocalIP(t *testing.T) {
	t.Run("resolveLocalIP() => localhost", func(t *testing.T) {
		r := resolveLocalIP(net.ParseIP("127.0.0.1"))
		if r.String() != "127.0.0.1" {
			t.Error("Failed to resolve localhost properly")
		}
	})

	t.Run("resolveLocalIP() => integration", func(t *testing.T) {
		host, ex := os.LookupEnv("SIP_HOST")

		if !def.Integration() || !ex {
			t.Skip()
		}

		r := resolveLocalIP(net.ParseIP(host))
		if r.String() == "127.0.0.1" {
			t.Error("Failed to resolve integration interface properly")
		}
	})
}
