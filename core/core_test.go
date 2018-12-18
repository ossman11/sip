package core

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/ossman11/sip/core/api"
	"github.com/ossman11/sip/core/test"
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

func ensureCrt() error {
	// Ensure that the certificate files are available
	os.MkdirAll("crt", os.ModePerm)

	if _, err := os.Stat("crt/server.crt"); os.IsNotExist(err) {
		_, err = copy("../crt/server.crt", "crt/server.crt")
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat("crt/server.key"); os.IsNotExist(err) {
		_, err = copy("../crt/server.key", "crt/server.key")
		if err != nil {
			return err
		}
	}
	return nil
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
		err := ensureCrt()
		if err != nil {
			t.Errorf("Failed to copy certificates, because: %v", err)
		}

		test.FindPort()
		defer test.OpenPort()

		res := NewServer()
		go res.Start()

		httpClient := test.HttpClient()

		for true {
			_, err := httpClient.Get(test.HttpServer())

			if err == nil {
				break
			}
		}
	})
}
