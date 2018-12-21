package index

import (
	"fmt"
	"net"
	"testing"

	"github.com/ossman11/sip/core/def"
)

func TestNetwork_GetRoutes(t *testing.T) {
	t.Run("GetRoutes() => ", func(t *testing.T) {
		network := Network{}

		n1 := Node{
			IP:   net.ParseIP("123.123.123.123"),
			ID:   "Node1",
			Type: Unknown,
			Port: def.GetPort(),
		}

		n2 := Node{
			IP:   net.ParseIP("123.123.123.123"),
			ID:   "Node2",
			Type: Unknown,
			Port: def.GetPort(),
		}

		c1 := NewConnection([]*Node{&n1, &n2})

		i1 := Index{}
		i1.Init()
		i1.Nodes[n1.ID] = &n1
		i1.Nodes[n2.ID] = &n2

		i1.Connections[c1.ID] = &c1

		i2 := Index{}
		i2.Init()
		i2.Nodes[n1.ID] = &n1
		i2.Nodes[n2.ID] = &n2

		i2.Connections[c1.ID] = &c1

		network.Add(&i1, n1.ID)
		network.Add(&i2, n2.ID)

		path := network.GetRoutes(n1.ID)
		fmt.Println(path)
	})
}

func TestNetwork_Path(t *testing.T) {
	t.Run("Path() => Direct Link", func(t *testing.T) {
		network := Network{}

		n1 := Node{
			IP:   net.ParseIP("123.123.123.123"),
			ID:   "Node1",
			Type: Unknown,
			Port: def.GetPort(),
		}

		n2 := Node{
			IP:   net.ParseIP("123.123.123.123"),
			ID:   "Node2",
			Type: Unknown,
			Port: def.GetPort(),
		}

		i1 := Index{}
		i1.Init()
		i1.Nodes[n1.ID] = &n1
		i1.Nodes[n2.ID] = &n2

		i2 := Index{}
		i2.Init()
		i2.Nodes[n1.ID] = &n1
		i2.Nodes[n2.ID] = &n2

		network.Add(&i1, n1.ID)
		network.Add(&i2, n2.ID)

		path := network.Path(n1.ID, n2.ID)
		fmt.Println(path)
	})
}
