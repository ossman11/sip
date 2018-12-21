package index

type Network struct {
	Indexs map[ID]*Index
	Routes map[ID][]*Route
}

func (n *Network) Add(i *Index, id ID) {
	if n.Indexs == nil {
		n.Indexs = map[ID]*Index{}
	}
	n.Indexs[id] = i
}

func (n *Network) Merge(a *Network) bool {
	ret := false

	for ck, cv := range a.Indexs {
		_, ex := n.Indexs[ck]
		if !ex {
			ret = true
			n.Add(cv, ck)
		}
	}

	return ret
}

func (n *Network) Has(id ID) bool {
	_, ex := n.Indexs[id]
	return ex
}

func (n *Network) GetRoutes(s ID) []*Route {

	r, ex := n.Routes[s]

	// TODO: ensure that no updates are required
	if ex {
		return r
	}

	i := n.Indexs[s]

	for ck, cv := range i.Connections {

		n0 := cv.Nodes[0]
		n1 := cv.Nodes[1]

		cr := n.GetRoutes()

		for _, crv := range cr {
			if !crv.Has(s) {
				nr := ExtendRoute(s, crv)
				r = append(r, &nr)
			}
		}
	}

	return r
}

func (n *Network) Path(s, t ID) []ID {
	p := make([]ID, 1)

	p[0] = t

	return p
}
