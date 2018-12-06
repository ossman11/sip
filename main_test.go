package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ossman11/sip/core/test"
)

func TestMain(m *testing.M) {
	flag.Parse()

	test.Integration()

	result := m.Run()

	os.Exit(result)
}

func TestCore(t *testing.T) {

	if !test.Integration() {
		t.Skip()
	}

	httpClient := test.HttpClient()
	res, err := httpClient.Get(test.HttpServer())

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
