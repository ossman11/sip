// +build !js

package index

import (
	"net"
)

func resolveLocalIP(r net.IP) net.IP {
	ret := "127.0.0.1"

	con, err := net.Dial("ip4:1", r.String())
	if err == nil {
		ret = con.LocalAddr().String()
		con.Close()
	}

	return net.ParseIP(ret)
}
