package def

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
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
	httpServer := "https://localhost:" + strconv.Itoa(GetPort())

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

func FindPort() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	os.Setenv("PORT", strconv.Itoa(port))
}
