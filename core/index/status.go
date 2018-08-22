package index

type Status int

func (s Status) String() string {
	if s < 0 {
		return "idle"
	}
	if s == Scanning {
		return "scanning"
	}
	if s == Indexing {
		return "indexing"
	}
	if s == Indexed {
		return "indexed"
	}
	return "unknown"
}

const (
	Idle     Status = -1
	Scanning Status = 1
	Indexing Status = 2
	Indexed  Status = 3
)
