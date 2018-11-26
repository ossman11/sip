package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/ossman11/sip/core/def"
)

var logSplitter = "====================================================="

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
	// Prepare local server url for easy consumption
	httpServer = "https://localhost:" + strconv.Itoa(def.GetPort())
)

// Define test parameters
var (
	integration = flag.Bool("integration", true, "Execute integration tests.")
)

func execCmd(str string) {
	src, err := filepath.Abs(ParseCommand(str))
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(src)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func setup() {
	fmt.Println(logSplitter)
	fmt.Println("Integration test setup: START")
	fmt.Println(logSplitter)
	// Ensure that older versions are cleaned up
	fmt.Println("Clean log:")
	execCmd("containers/kube/clean")

	// Ensure that the latest changes are build
	fmt.Println("Build log:")
	execCmd("containers/build")

	// Ensure that the latest changes are deployed
	fmt.Println("Deploy log:")
	execCmd("containers/kube/deploy")

	// Ensure that the latest changes are exposed
	fmt.Println("Await log:")
	execCmd("containers/kube/await")

	// Ensure that localhost host a server instance
	fmt.Println("Local server log:")
	go main()

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
			os.Exit(1)
		}
	}

	fmt.Println(logSplitter)
	fmt.Println("Integration test setup: DONE")
	fmt.Println(logSplitter)
}

func teardown() {
	fmt.Println(logSplitter)
	fmt.Println("Integration test teardown: START")
	fmt.Println(logSplitter)
	// Ensure that the created test resources are cleaned
	fmt.Println("Clean log:")
	execCmd("containers/kube/clean")

	fmt.Println(logSplitter)
	fmt.Println("Integration test teardown: DONE")
	fmt.Println(logSplitter)
}

func TestMain(m *testing.M) {
	flag.Parse()

	if *integration {
		setup()
	}

	result := m.Run()

	if *integration {
		teardown()
	}

	os.Exit(result)
}

func TestCore(t *testing.T) {

	if !*integration {
		t.Skip()
	}

	res, err := httpClient.Get(httpServer)

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
