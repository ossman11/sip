package def

const (
	// Port is the default port on which sip nodes run
	Port int = 1670
)

const (
	// APIIndex provides the root url for the index API
	APIIndex string = "/index"
	// APIIndexJoin provides the join url for the index API
	APIIndexJoin string = APIIndex + "/join"
	// APIIndexNodes provides the nodes url for the index API
	APIIndexNodes string = APIIndex + "/nodes"
)
