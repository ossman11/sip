package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ossman11/sip/core/def"
)

func TestMain(m *testing.M) {
	def.Integration()
	result := m.Run()
	os.Exit(result)
}

func TestCore(t *testing.T) {
	if !def.Integration() {
		t.Skip()
	}

	def.FindPort()

	go main()

	httpClient := def.HttpClient()
	httpServer := def.HttpServer()

	attempts := 0
	for true {
		res, err := httpClient.Get(httpServer)

		if err == nil {
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)

			if err != nil {
				t.Errorf("Failed to read response body with the error: %s", err)
			}

			fmt.Print(string(body))
			return
		}

		attempts++
		if attempts > 20 {
			t.Errorf("Failed to connect to the local server with the error: %s", err)
		}
	}
}
