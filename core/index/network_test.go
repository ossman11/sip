package index

import (
	"fmt"
	"math"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/ossman11/sip/core/def"
)

var testRandom = rand.New(rand.NewSource(time.Now().UnixNano()))

func mockNode(id ID) *Node {
	n := Node{
		IP:   net.ParseIP("123.123.123.123"),
		ID:   id,
		Type: Unknown,
		Port: def.GetPort(),
	}
	return &n
}

func mockConn(n1, n2 *Node) *Connection {
	c := NewConnection([]*Node{n1, n2})
	return &c
}

func mockIndex() *Index {
	i := Index{}
	i.Init()
	return &i
}

func addConn(i *Index, c *Connection) {
	i.Connections[c.ID] = c
	i.Nodes[c.Nodes[0].ID] = c.Nodes[0]
	i.Nodes[c.Nodes[1].ID] = c.Nodes[1]
}

func testNetworkDirect() *Network {
	network := Network{}

	n1 := mockNode("Node1")
	n2 := mockNode("Node2")

	c1 := mockConn(n1, n2)

	i1 := mockIndex()
	addConn(i1, c1)

	i2 := mockIndex()
	addConn(i2, c1)

	network.Add(i1, n1.ID)
	network.Add(i2, n2.ID)

	return &network
}

func testNetworkNested() *Network {
	network := Network{}

	n1 := mockNode("Node1")
	n2 := mockNode("Node2")
	n3 := mockNode("Node3")

	c1 := mockConn(n1, n2)
	c2 := mockConn(n2, n3)

	i1 := mockIndex()
	addConn(i1, c1)

	i2 := mockIndex()
	addConn(i2, c1)
	addConn(i2, c2)

	i3 := mockIndex()
	addConn(i3, c2)

	network.Add(i1, n1.ID)
	network.Add(i2, n2.ID)
	network.Add(i3, n3.ID)

	return &network
}

func testNetworkCyclic() *Network {
	network := Network{}

	n1 := mockNode("Node1")
	n2 := mockNode("Node2")
	n3 := mockNode("Node3")

	c1 := mockConn(n1, n2)
	c2 := mockConn(n2, n3)
	c3 := mockConn(n1, n3)

	i1 := mockIndex()
	addConn(i1, c1)
	addConn(i1, c3)

	i2 := mockIndex()
	addConn(i2, c1)
	addConn(i2, c2)

	i3 := mockIndex()
	addConn(i3, c2)
	addConn(i3, c3)

	network.Add(i1, n1.ID)
	network.Add(i2, n2.ID)
	network.Add(i3, n3.ID)

	return &network
}

func testNetworkDeep() *Network {
	network := Network{}

	i := 0
	for i < 100 {
		i++
		cid := ID("Node" + strconv.Itoa(i))
		cn := mockNode(cid)
		ci := mockIndex()
		ci.Nodes[cn.ID] = cn
		network.Add(ci, cn.ID)
	}

	i = 0
	for i < 100 {
		i++
		cid := ID("Node" + strconv.Itoa(i))
		ci := network.Indexs[cid]

		ccnr := math.Ceil(testRandom.Float64()*20) + 2

		for ccnr > 0 {
			tid := cid
			for tid == cid {
				tid = ID("Node" + strconv.Itoa(int(math.Ceil(testRandom.Float64()*100))))
			}

			ti := network.Indexs[tid]
			cc := mockConn(ci.Nodes[cid], ti.Nodes[tid])

			addConn(ci, cc)
			addConn(ti, cc)

			ccnr--
		}
	}

	return &network
}

func TestNetwork_GetRoutes(t *testing.T) {
	t.Run("GetRoutes() => direct", func(t *testing.T) {
		network := testNetworkDirect()
		fmt.Println(network.Paths)
	})

	t.Run("GetRoutes() => nested", func(t *testing.T) {
		network := testNetworkNested()
		fmt.Println(network.Paths)
	})

	t.Run("GetRoutes() => cyclic", func(t *testing.T) {
		network := testNetworkCyclic()
		fmt.Println(network.Paths)
	})

	t.Run("GetRoutes() => deep", func(t *testing.T) {
		network := testNetworkDeep()
		err, p := network.Path("Node1", "Node100")
		fmt.Println(err)
		fmt.Println(p)
	})
}

func TestNetwork_Path(t *testing.T) {
	t.Run("Path() => Direct Link", func(t *testing.T) {
		network := testNetworkDirect()

		err, p := network.Path("Node1", "Node2")
		fmt.Println(err)
		fmt.Println(p)
	})
}
