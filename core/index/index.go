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
	Connections map[ID]Connection
	Nodes       map[ID]Node
	scanner     Scan
	httpClient  *http.Client
	updateChan  chan bool
	updateNext  bool
}

func (i *Index) Add(c Node) {
	i.Nodes[c.ID] = c
	i.JoinNode(c)
}

func (i *Index) AddAll(n map[ID]Node) {
	for k := range n {
		i.Add(n[k])
	}
}

func (i *Index) AddCon(c Connection) {
	i.Connections[c.ID] = c
}

func (i *Index) JoinNode(newNode Node) {
	node := ThisNode(i, newNode.IP)
	newConNodes := []*Node{
		&node,
		&newNode,
	}
	newCon := NewConnection(newConNodes)

	if newNode.IP.String() != node.IP.String() && newNode.ID != node.ID {
		_, nodeExists := i.Nodes[newNode.ID]
		_, conExists := i.Connections[newCon.ID]
		if !nodeExists || !conExists {
			if !conExists {
				i.AddCon(newCon)
			}

			if !nodeExists {
				i.Add(newNode)
			}

			i.Join(newNode.IP)
		}
	}
}

func (i *Index) Join(ip net.IP) {
	str := ip.String() + ":" + strconv.Itoa(def.Port)
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
			fmt.Println("JOIN: ", node.IP.String(), " -> ", ip.String())

			newNode := Node{}
			dec := json.NewDecoder(res.Body)
			err = dec.Decode(&newNode)

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
			_, ex := i.Nodes[newNode.ID]
			change = ex || change
			if !ex {
				i.Add(newNode)
			}

			_, ex = i.Connections[newCon.ID]
			change = ex || change
			if !ex {
				i.AddCon(newCon)
			}

			if change {
				go i.Update()
			}
		}
	}
}

func (i *Index) Merge(n *Index) bool {
	ret := false

	for ck, cv := range n.Connections {
		_, ex := i.Connections[ck]
		if !ex {
			ret = true
			i.AddCon(cv)
		}
	}

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
			"https://"+v.IP.String()+":"+strconv.Itoa(def.Port)+def.APIIndex,
			"application/json",
			bytes.NewBuffer(bod))

		if err == nil {
			fmt.Println("INDEX: ", thisNode.IP.String(), " -> ", v.IP.String())
			resBod := Index{}
			dec := json.NewDecoder(res.Body)
			err = dec.Decode(&resBod)

			i.Merge(&resBod)
		}
	}

	ch := i.updateChan
	i.updateChan = nil
	close(ch)
}

func (i *Index) Scan() {
	i.Init()
	i.Status = Scanning
	i.scanner.Scan()
	i.Status = Indexing
	i.Update()
	i.Status = Idle
}

func (i *Index) Init() {
	if i.Nodes == nil {
		i.Nodes = map[ID]Node{}
	}

	if i.Connections == nil {
		i.Connections = map[ID]Connection{}
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

	i.Add(ThisNode(i, net.ParseIP("127.0.0.1")))
	if i.scanner == (Scan{}) {
		i.scanner = NewScan(i)
	}
	if i.Status == 0 {
		i.Status = Idle
	}
}
