package index

type Network struct {
	Indexs map[ID]*Index
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
