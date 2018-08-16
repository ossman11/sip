package index

import (
	"fmt"
	"net"
)

// TODO: aggregate whole network into nodes
// TODO: allow for redirect indexing nodes
type Index struct {
	Type        Type
	Connections map[string]Connection
	Nodes       map[string]Connection
	scanner     Scan
}

func (i *Index) Add(id string, c Connection) {
	if i.Connections == nil {
		i.Connections = map[string]Connection{}
	}
	i.Connections[id] = c
}

func (i *Index) AddAll(n map[string]Connection) {
	for k := range n {
		i.Add(k, n[k])
	}
}

func (i *Index) Get(id string, ip net.IP) ID {
	ID := i.Connections[id].ID
	curID := &ID
	return curID.Out(ip)
}

func (i *Index) GetAll(ip net.IP) map[string]ID {
	r := map[string]ID{}
	for k := range i.Connections {
		c := i.Get(k, ip)
		r[c.String()] = c
	}
	return r
}

func (i *Index) Scan() {
	i.Init()
	i.scanner.Scan()
	i.AddAll(i.scanner.Result)
	fmt.Println(i.scanner.Result)
}

func (i *Index) Status() bool {
	i.Init()
	return i.scanner.Running
}

func (i *Index) Init() {
	if i.scanner.Result == nil {
		i.scanner = NewScan()
	}
}
