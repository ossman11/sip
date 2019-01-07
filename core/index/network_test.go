package index

import (
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

	nodeCount := 100
	maxConnections := 25

	i := 0
	for i < nodeCount {
		i++
		cid := ID("Node" + strconv.Itoa(i))
		cn := mockNode(cid)
		ci := mockIndex()
		ci.Nodes[cn.ID] = cn
		network.Add(ci, cn.ID)
	}

	i = 0
	for i < nodeCount {
		i++
		cid := ID("Node" + strconv.Itoa(i))
		ci := network.Indexs[cid]

		ccnr := math.Ceil(testRandom.Float64()*float64(maxConnections)) + 2

		for ccnr > 0 {
			tid := cid
			for tid == cid {
				tid = ID("Node" + strconv.Itoa(int(math.Ceil(testRandom.Float64()*float64(nodeCount)))))
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

func TestNetwork_Path(t *testing.T) {
	t.Run("Path() => direct", func(t *testing.T) {
		network := testNetworkDirect()
		err, p := network.Path("Node1", "Node2")

		if err != nil {
			t.Error(err)
		}

		if p[1] == nil {
			t.Error("Failed to find the direct link inside the direct network.\n", p)
		}
	})

	t.Run("Path() => nested", func(t *testing.T) {
		network := testNetworkNested()
		err, p := network.Path("Node1", "Node3")

		if err != nil {
			t.Error(err)
		}

		if p[1] != nil {
			t.Error("Found a direct link inside the nested network.\n", p)
		}

		if p[2] == nil {
			t.Error("Failed to find the nested link inside the nested network.\n", p)
		}
	})

	t.Run("Path() => cyclic", func(t *testing.T) {
		network := testNetworkCyclic()
		err, p := network.Path("Node1", "Node3")

		if err != nil {
			t.Error(err)
		}

		if p[1] == nil {
			t.Error("Failed to find the direct link inside the cyclic network.\n", p)
		}

		if p[2] == nil {
			t.Error("Failed to find the nested link inside the cyclic network.\n", p)
		}
	})

	t.Run("Path() => deep", func(t *testing.T) {
		network := testNetworkDeep()
		err, p := network.Path("Node1", "Node100")

		if err != nil {
			t.Error(err)
		}

		for _, sp := range p {
			for _, cv := range sp {
				last := ID("Node1")
				for i, n := range cv.Nodes {
					if i == 0 {
						if last == n {
							continue
						}
						t.Error("Failed to traverse path as it does not start at the correct Node.")
					}

					if last == n {
						t.Error("Failed to traverse path as path contains same node ID twice.")
					}

					cn := network.Indexs[n]
					if cn.Connections[NewIDConnection([]*Node{cn.Nodes[n], network.Indexs[last].Nodes[last]})] == nil {
						t.Error("Failed to traverse path as it is using invalid connections.\n", cv.Nodes)
					}

					last = n
				}
			}
		}
	})

	t.Run("Path() => missing path", func(t *testing.T) {
		network := testNetworkDirect()

		n3 := mockNode("Node3")
		i3 := mockIndex()
		network.Add(i3, n3.ID)

		err, p := network.Path("Node1", "Node3")

		if p != nil {
			t.Error("Succeeded to find path that should not exist.")
		}

		if err == nil {
			t.Error("Did not retrieve error message from failing Path lookup.")
		}
	})

	t.Run("Path() => missing target index", func(t *testing.T) {
		network := testNetworkDirect()
		network.Paths = nil
		err, p := network.Path("Node1", "Node3")

		if p != nil {
			t.Error("Succeeded to find path that should not exist.")
		}

		if err == nil {
			t.Error("Did not retrieve error message from failing Path lookup.")
		}
	})

	t.Run("Path() => missing start index", func(t *testing.T) {
		network := testNetworkDirect()
		network.Paths = nil
		err, p := network.Path("Node3", "Node1")

		if p != nil {
			t.Error("Succeeded to find path that should not exist.")
		}

		if err == nil {
			t.Error("Did not retrieve error message from failing Path lookup.")
		}
	})
}

func TestNetwork_AddRoute(t *testing.T) {
	t.Run("AddRoute() => initialize", func(t *testing.T) {
		n1 := mockNode("Node1")
		n2 := mockNode("Node2")

		r := NewRoute([]ID{n1.ID, n2.ID})

		network := Network{}
		network.AddRoute(r.ID, &r)

		if network.Paths == nil || network.Paths[r.ID] == nil {
			t.Error("Failed to setup Paths cache properly for new Route.")
		}
	})
}
