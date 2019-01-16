// +build js

package index

import (
	"net"
)

func resolveLocalIP(r net.IP) net.IP {
	return net.ParseIP("0.0.0.0")
}
