package index

type Route struct {
	Nodes []ID
}

func (r *Route) Has(id ID) bool {
	for _, cv := range r.Nodes {
		if cv == id {
			return true
		}
	}
	return false
}

func NewRoute(n []ID) Route {
	return Route{
		Nodes: n,
	}
}

func ExtendRoute(s ID, r *Route) Route {
	return Route{
		Nodes: append([]ID{s}, r.Nodes...),
	}
}
