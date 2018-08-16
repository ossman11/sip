package index

import "net"

// TODO: alive check
type Connection struct {
	IP   net.IP
	ID   ID
	Type Type
}
