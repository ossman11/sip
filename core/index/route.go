package index

type Route struct {
	Nodes []ID
	ID    ID
}

func (r *Route) String() string {
	ret := ""
	for _, cv := range r.Nodes {
		if ret != "" {
			ret += ","
		}
		ret += string(cv)
	}
	return ret
}
func (r *Route) Has(id ID) bool {
	for _, cv := range r.Nodes {
		if cv == id {
			return true
		}
	}
	return false
}

func (r *Route) Equal(t *Route) bool {
	if r.Len() != t.Len() {
		return false
	}

	for i := range r.Nodes {
		if r.Nodes[i] != t.Nodes[i] {
			return false
		}
	}

	return true
}

func (r *Route) Next() Route {
	return NewRoute(r.Nodes[1:])
}

func (r *Route) Len() int {
	return len(r.Nodes) - 1
}

func NewRoute(n []ID) Route {
	return Route{
		Nodes: n,
		ID:    NewIDRoute(n),
	}
}

func ExtendRoute(s ID, r *Route) Route {
	return NewRoute(append([]ID{s}, r.Nodes...))
}
