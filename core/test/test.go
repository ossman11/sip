package test

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ossman11/sip/core/def"
)

var (
	logSplitter = "====================================================="
	ready       = false
	integration = flag.Bool("integration", false, "Execute integration tests.")
)

func HttpClient() *http.Client {
	// Disable tls until properly adopted
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// Expose insecure http client
	httpClient := &http.Client{
		Transport: tr,
	}

	return httpClient
}

func HttpServer() string {
	// Prepare local server url for easy consumption
	httpServer := "https://localhost:" + strconv.Itoa(def.GetPort())

	return httpServer
}

func Integration() bool {
	flag.Parse()

	// Check if integration is enabled and if not return without doing anything
	if !*integration {
		return false
	}

	if *integration && ready {
		return true
	}

	httpClient := HttpClient()
	httpServer := HttpServer()

	attempts := 0
	for true {
		_, err := httpClient.Get(httpServer)

		if err == nil {
			break
		}

		fmt.Printf(".")

		attempts++
		if attempts > 20 {
			fmt.Printf("Failed to connect to the local server with the error: %s\n", err)
			tf := false
			integration = &tf
			return false
		}
	}

	fmt.Println(logSplitter)
	fmt.Println("Integration test setup: DONE")
	fmt.Println(logSplitter)

	ready = true

	return true
}

var lastPort = 0

func FindPort() {
	port := def.GetPort()
	lastPort = port

	for true {
		_, err := HttpClient().Get(HttpServer())

		if err == nil {
			port++
			os.Setenv("PORT", strconv.Itoa(port))
		} else {
			break
		}
	}
}

func OpenPort() {
	os.Setenv("PORT", strconv.Itoa(lastPort))
}
