package core

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/ossman11/sip/core/api"
	"github.com/ossman11/sip/core/def"
)

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}

	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
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
		err := GenCrt()
		if err != nil {
			t.Errorf("Failed to generate certificates, because: %v", err)
		}

		def.FindPort()

		res := NewServer()
		go res.Start()

		httpClient := def.HttpClient()

		for true {
			_, err := httpClient.Get(def.HttpServer())

			if err == nil {
				break
			}
		}
	})
}
