package def

import (
	"os"
	"strconv"
)

const (
	// Port is the default port on which sip nodes run
	Port int = 1670
)

var curPort int

func GetPort() int {
	if curPort != 0 {
		return curPort
	}

	curPort = Port
	envPort, ex := os.LookupEnv("PORT")
	if ex {
		tmpPort, err := strconv.Atoi(envPort)
		if err != nil {
			tmpPort = Port
		}
		curPort = tmpPort
	}
	return GetPort()
}

const (
	// APIIndex provides the root url for the index API
	APIIndex string = "/index"
	// APIIndexJoin provides the join url for the index API
	APIIndexJoin string = APIIndex + "/join"
	// APIIndexCollect provides the collect url for the index API
	APIIndexCollect string = APIIndex + "/collect"
)
