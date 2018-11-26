package index

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/ossman11/sip/core/def"
)

// TODO: aggregate whole network into nodes
// TODO: allow for redirect indexing nodes
type Index struct {
	Type        Type
	Status      Status
	Connections map[ID]*Connection
	Nodes       map[ID]*Node
	scanner     Scan
	httpClient  *http.Client
	updateChan  chan bool
	updateNext  bool
}

func (i *Index) Add(c *Node) {
	i.Nodes[c.ID] = c
	i.JoinNode(c)
}

func (i *Index) AddAll(n map[ID]*Node) {
	for k := range n {
		i.Add(n[k])
	}
}

func (i *Index) AddCon(c *Connection) {
	i.Connections[c.ID] = c
}

func (i *Index) JoinNode(newNode *Node) {
	node := ThisNode(i, newNode.IP)

	if node.IP.String() == "127.0.0.1" {
		return
	}

	if newNode.IP.String() != node.IP.String() && newNode.ID != node.ID {
		_, nodeExists := i.Nodes[newNode.ID]
		if !nodeExists  {
			i.Join(newNode.IP, newNode.Port)
		}
	}
}

func (i *Index) Join(ip net.IP, port int) {
	if port == 0 {
		port = def.Port
	}
	str := ip.String() + ":" + strconv.Itoa(port)
	node := ThisNode(i, ip)
	bod, err := json.Marshal(node)
	if err != nil {
		fmt.Println(err)
		return
	}

	if ip.String() != node.IP.String() {
		res, err := i.httpClient.Post(
			"https://"+str+def.APIIndexJoin,
			"application/json",
			bytes.NewBuffer(bod))

		if err == nil {
			// fmt.Println("JOIN: ", node.IP.String(), " -> ", ip.String())

			newNode := Node{}
			dec := json.NewDecoder(res.Body)
			err = dec.Decode(&newNode)
			newNode.IP = ip
			newNode.Port = port

			if err != nil {
				fmt.Println(err)
				return
			}

			newConNodes := []*Node{
				&node,
				&newNode,
			}
			newCon := NewConnection(newConNodes)

			change := false
			exNode, ex := i.Nodes[newNode.ID]
			change = ex || change
			updateNode := ex && (exNode.IP.String() != newNode.IP.String() || exNode.Port != newNode.Port)
			if !ex || updateNode {
				i.Add(&newNode)
			}

			_, ex = i.Connections[newCon.ID]
			change = ex || change
			if !ex || updateNode {
				i.AddCon(&newCon)
			}

			if change {
				go i.Update()
			}
		}
	}
}

func (i *Index) Merge(n *Index) bool {
	ret := false

	/*
		for ck, cv := range n.Connections {
			_, ex := i.Connections[ck]
			if !ex {
				ret = true
				i.AddCon(cv)
			}
		}
	*/

	for nk, nv := range n.Nodes {
		_, ex := i.Nodes[nk]
		if !ex {
			ret = true
			i.Add(nv)
		}
	}

	return ret
}

func (i *Index) Update() {
	// If already updating update again later
	if i.updateChan != nil {
		// If alreadt waiting to update again continue
		if i.updateNext {
			return
		}
		i.updateNext = true
		<-i.updateChan
		i.updateNext = false
	}

	i.updateChan = make(chan bool)

	bod, _ := json.Marshal(i)
	for _, v := range i.Nodes {
		thisNode := ThisNode(i, v.IP)
		// Skip updating self
		if v.ID == thisNode.ID {
			continue
		}
		res, err := i.httpClient.Post(
			"https://"+v.IP.String()+":"+strconv.Itoa(v.Port)+def.APIIndex,
			"application/json",
			bytes.NewBuffer(bod))

		if err == nil {
			// fmt.Println("INDEX: ", thisNode.IP.String(), " -> ", v.IP.String())
			resBod := Index{}
			dec := json.NewDecoder(res.Body)
			err = dec.Decode(&resBod)

			i.Merge(&resBod)
		}
	}

	ch := i.updateChan
	i.updateChan = nil
	if ch != nil {
		close(ch)
	}
}

func (i *Index) Scan() {
	i.Init()
	i.Status = Scanning
	i.scanner.Scan()
	i.Status = Indexing
	i.Update()
	i.Status = Idle
}

func (i *Index) Collect(n *Network) {
	thisNode := ThisNode(i, net.ParseIP("127.0.0.1"))
	n.Add(i, thisNode.ID)

	bod, _ := json.Marshal(n)
	for _, v := range i.Nodes {
		if n.Has(v.ID) {
			continue
		}

		res, err := i.httpClient.Post(
			"https://"+v.IP.String()+":"+strconv.Itoa(v.Port)+def.APIIndexCollect,
			"application/json",
			bytes.NewBuffer(bod))

		if err == nil {
			resBod := Network{}
			dec := json.NewDecoder(res.Body)
			err = dec.Decode(&resBod)

			n.Merge(&resBod)
		}
	}
}

func (i *Index) Init() {
	if i.Nodes == nil {
		i.Nodes = map[ID]*Node{}
	}

	if i.Connections == nil {
		i.Connections = map[ID]*Connection{}
	}

	timeout := time.Duration(30 * time.Second)

	// Always scan without security enabled
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	i.httpClient = &http.Client{
		Timeout:   timeout,
		Transport: tr,
	}

	curNode := ThisNode(i, net.ParseIP("127.0.0.1"))

	i.Add(&curNode)
	if i.scanner == (Scan{}) {
		i.scanner = NewScan(i)
	}
	if i.Status == 0 {
		i.Status = Idle
	}
}
