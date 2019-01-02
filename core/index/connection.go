package index

type Connection struct {
	ID    ID
	Nodes []*Node
}

func NewConnection(n []*Node) Connection {
	return Connection{
		ID:    NewIDConnection(n),
		Nodes: n,
	}
}
