package def

import (
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// Port is the default port on which sip nodes run
	Port int = 1670
)

func GetParent() (string, int) {
	fileName := os.Args[0]
	fileName = strings.Trim(
		strings.Replace(
			fileName,
			filepath.Dir(fileName),
			"",
			1,
		),
		"\\/.exe",
	)

	if host, ps, err := net.SplitHostPort(strings.Replace(fileName, "_", ":", 1)); err == nil {
		p, err := strconv.Atoi(ps)
		if err == nil {
			return host, p
		}
	}
	return "", 0
}

func GetPort() int {
	envPort, ex := os.LookupEnv("PORT")
	if ex {
		tmpPort, err := strconv.Atoi(envPort)
		if err == nil {
			return tmpPort
		}
	}

	ph, pp := GetParent()
	if ph != "" {
		return pp
	}
	return Port
}

const (
	// APIIndex provides the root url for the index API
	APIIndex string = "/index"
	// APIIndexJoin provides the join url for the index API
	APIIndexJoin string = APIIndex + "/join"
	// APIIndexCollect provides the collect url for the index API
	APIIndexCollect string = APIIndex + "/collect"
	// APIIndexCall provides the call url for the index API
	APIIndexCall string = APIIndex + "/call"
)
