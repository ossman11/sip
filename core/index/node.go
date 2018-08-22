package index

import (
	"net"
)

// TODO: alive check
type Node struct {
	IP   net.IP
	ID   ID
	Type Type
}

func resolveLocalIP(r net.IP) net.IP {
	faces, err := net.Interfaces()
	if err == nil {
		for _, v := range faces {
			as, err := v.Addrs()
			if err != nil {
				continue
			}
			for _, a := range as {
				ipnet, ok := a.(*net.IPNet)
				if !ok {
					continue
				}

				ip4 := ipnet.IP.To4()
				rip4 := r.To4()
				if ip4 != nil && rip4 != nil {
					mask := ipnet.Mask
					match := true
					for i := range mask {
						if mask[i] == 255 && ip4[i] != rip4[i] {
							match = false
						}
					}
					if match {
						return ipnet.IP
					}
				}
			}
		}
	}
	return net.ParseIP("127.0.0.1")
}
func ThisNode(i *Index, face net.IP) Node {
	t := Unknown
	if i != nil {
		t = i.Type
	}
	return Node{
		IP:   resolveLocalIP(face),
		ID:   ParseHWID(),
		Type: t,
	}
}
