package index

import (
	"net"

	"github.com/ossman11/sip/core/def"
)

// TODO: alive check
type Node struct {
	IP   net.IP
	ID   ID
	Type Type
	Port int
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
		Port: def.GetPort(),
	}
}
