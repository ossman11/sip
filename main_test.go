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

	httpClient := def.HttpClient()
	res, err := httpClient.Get(def.HttpServer())

	if err != nil {
		t.Errorf("Failed to connect to the local server with the error: %s", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Errorf("Failed to read response body with the error: %s", err)
	}

	fmt.Print(string(body))
}
