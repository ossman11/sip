package core

import (
	"crypto/tls"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/ossman11/sip/core/api"
	"github.com/ossman11/sip/core/def"
)

// Define http information
var (
	// Disable tls until properly adopted
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// Expose insecure http client
	httpClient = &http.Client{
		Transport: tr,
	}
)

func getLocalServer() string {
	return "https://localhost:" + strconv.Itoa(def.GetPort())
}

func TestNewServer(t *testing.T) {
	t.Run("NewServer() => default", func(t *testing.T) {
		res := NewServer()
		if !res.ready {
			t.Errorf("Failed to use NewServer(), because ready was %v.", res.ready)
		}

		if res.busy == nil {
			t.Errorf("Failed to use NewServer(), because busy was %v.", res.busy)
		}

		if res.server == nil {
			t.Errorf("Failed to use NewServer(), because server was %v.", res.server)
		}

		if res.handler == nil {
			t.Errorf("Failed to use NewServer(), because handler was %v.", res.handler)
		}
	})

	t.Run("NewServer() => secondary Init()", func(t *testing.T) {
		res := NewServer()
		res.Init()
	})

	t.Run("NewServer() => AddApis()", func(t *testing.T) {
		res := NewServer()
		res.AddApis(api.Empty{})
	})

	t.Run("NewServer() => Start()", func(t *testing.T) {
		res := NewServer()

		port := def.GetPort()

		for true {
			_, err := httpClient.Get(getLocalServer())

			if err == nil {
				port++
				os.Setenv("PORT", strconv.Itoa(port))
			} else {
				break
			}
		}

		go res.Start()

		for true {
			_, err := httpClient.Get(getLocalServer())

			if err == nil {
				break
			}
		}
	})
}