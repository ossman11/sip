package index

type Type int

func (t Type) String() string {
	if t < 0 {
		return "unknown"
	}
	if t == Local {
		return "local"
	}
	if t == Redirect {
		return "redirect"
	}
	return "unknown"
}

const (
	Unknown  Type = -1
	Local    Type = 0
	Redirect Type = 1
)
