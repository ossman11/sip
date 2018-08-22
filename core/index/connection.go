package index

import (
	"crypto/sha256"
	"sort"
)

type Connection struct {
	ID    ID
	Nodes []*Node
}

func NewConnection(n []*Node) Connection {
	sort.Slice(n, func(i, j int) bool {
		return n[i].ID < n[j].ID
	})

	conID := ""

	for _, v := range n {
		conID += string(v.ID) + ":" + v.IP.String() + "\n"
	}

	b := sha256.Sum256([]byte(conID))
	id := ParseByte(b[:])
	return Connection{
		ID:    id,
		Nodes: n,
	}
}
