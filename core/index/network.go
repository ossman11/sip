package index

import (
	"errors"
	"sort"
)

var pathRedundancy = 8

type Network struct {
	Indexs map[ID]*Index
	Paths  map[ID]map[ID]map[int][]*Route
}

func (n *Network) Add(i *Index, id ID) {
	if n.Indexs == nil {
		n.Indexs = map[ID]*Index{}
	}
	n.Indexs[id] = i
	n.AddIndex(i, id, nil)
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

func (n *Network) AddIndex(i *Index, id ID, hit map[ID]bool) {
	// Ensure that the hit array is initialized
	if hit == nil {
		hit = map[ID]bool{}
	}

	hit[id] = true

	// Setup root Paths with index ID mapping
	if n.Paths == nil {
		n.Paths = map[ID]map[ID]map[int][]*Route{}
	}

	// Setup index Paths with route ID mapping
	if n.Paths[id] == nil {
		n.Paths[id] = map[ID]map[int][]*Route{}
	}

	for _, cv := range i.Connections {
		n0 := cv.Nodes[0]
		n1 := cv.Nodes[1]
		tn := n0

		if n0.ID == id {
			tn = n1
		}

		ti := n.Indexs[tn.ID]

		nr := NewRoute([]ID{id, tn.ID})
		n.AddRoute(id, &nr)

		// Target index is unknown so skipping extending paths
		if ti == nil {
			continue
		}

		if n.Paths[tn.ID] == nil || !hit[tn.ID] {
			n.AddIndex(ti, tn.ID, hit)
		}

		tp := n.Paths[tn.ID]

		for _, ccv := range tp {

			pr := make([]int, 0)
			for prk, prv := range ccv {
				if len(prv) > 0 {
					pr = append(pr, prk)
				}
			}
			sort.Ints(pr)

			// Take top most optimal routes only
			concount := 0
			for _, prk := range pr {
				if concount >= pathRedundancy {
					break
				}
				cccv := ccv[prk]
				for _, ccccv := range cccv {
					if concount >= pathRedundancy {
						break
					}
					if ccccv.Has(id) {
						continue
					}
					cnr := ExtendRoute(id, ccccv)
					n.AddRoute(id, &cnr)
					concount++
				}
			}
		}
	}
}

func (n *Network) AddRoute(i ID, r *Route) {
	// Setup root Paths with index ID mapping
	if n.Paths == nil {
		n.Paths = map[ID]map[ID]map[int][]*Route{}
	}

	// Setup index Paths with route ID mapping
	if n.Paths[i] == nil {
		n.Paths[i] = map[ID]map[int][]*Route{}
	}

	p := n.Paths[i]

	if p[r.ID] == nil {
		p[r.ID] = map[int][]*Route{}
	}

	ps := p[r.ID]

	if r.Len() > 1 {
		pr := make([]int, 0)
		for prk, prv := range ps {
			if len(prv) > 0 {
				pr = append(pr, prk)
			}
		}
		sort.Ints(pr)

		concount := 0
		addNew := true
		for i := range pr {
			concount = concount + len(ps[pr[i]])
			if concount >= pathRedundancy {
				delete(ps, pr[i])
				if addNew && r.Len() > pr[i] {
					addNew = false
				}
			}
		}
		if !addNew {
			return
		}
	}

	if ps[r.Len()] == nil {
		ps[r.Len()] = []*Route{}
	}

	rc := ps[r.Len()]
	has := false

	for _, cv := range rc {
		if r.Equal(cv) {
			has = true
		}
	}

	if !has {
		ps[r.Len()] = append(rc, r)
	}
}

func (n *Network) Path(s, t ID) (error, map[int][]*Route) {
	tr := NewRoute([]ID{s, t})

	if n.Paths == nil || n.Paths[s] == nil || n.Paths[s][tr.ID] == nil {
		if n.Indexs[s] == nil {
			errMsg := string("Failed to lookup path between Node \"" + s + "\" and \"" + t + "\"" +
				", because starting Node with ID: " + s + " was not found.")
			return errors.New(errMsg), nil
		}

		n.AddIndex(n.Indexs[s], s, nil)
	}

	if n.Paths[s][tr.ID] == nil {
		errMsg := string("Failed to lookup path between Node \"" + s + "\" and \"" + t + "\"" +
			", because no known path was found.")
		return errors.New(errMsg), nil
	}

	return nil, n.Paths[s][tr.ID]
}
