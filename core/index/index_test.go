package index

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ossman11/sip/core/test"
)

func TestMain(m *testing.M) {
	test.Integration()

	result := m.Run()

	os.Exit(result)
}

func newMockIndex() *Index {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))

	index := Index{}
	index.Init()
	index.httpClient = ts.Client()
	index.httpClient.Timeout = 100 * time.Nanosecond

	return &index
}

func TestIndex_Init(t *testing.T) {
	t.Run("Init()", func(t *testing.T) {
		index := Index{}
		index.Init()
		// TODO: Add checks
	})
}

func TestIndex_JoinNode(t *testing.T) {
	t.Run("JoinNode() => remote node", func(t *testing.T) {
		index := Index{}
		index.Init()
		index.httpClient.Timeout = 100 * time.Nanosecond

		newNode := Node{
			IP:   net.ParseIP("123.123.123.123"),
			ID:   "MOCKID",
			Type: Unknown,
			Port: 123,
		}

		index.JoinNode(&newNode)
	})
}

func TestIndex_Join(t *testing.T) {
	t.Run("Join() => Non existing node", func(t *testing.T) {
		index := Index{}
		index.Init()
		index.httpClient.Timeout = 100 * time.Nanosecond

		index.Join(net.ParseIP("123.123.123.123"), 123)
	})

	t.Run("Join() => Non existing node without port", func(t *testing.T) {
		index := Index{}
		index.Init()
		index.httpClient.Timeout = 100 * time.Nanosecond

		index.Join(net.ParseIP("123.123.123.123"), 0)
	})

	t.Run("Join() => Existing node", func(t *testing.T) {
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `{
				"IP":"127.0.0.1",
				"ID":"ABC",
				"Type":-1,
				"Port": 0
			}`)
		}))

		index := Index{}
		index.Init()

		index.httpClient = ts.Client()

		urlSplit := strings.SplitN(ts.URL, "://", 2)

		host, portStr, _ := net.SplitHostPort(urlSplit[1])
		port, _ := strconv.Atoi(portStr)

		index.Join(net.ParseIP(host), port)
	})

	t.Run("Join() => Existing node(invalid response)", func(t *testing.T) {
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `not a json`)
		}))

		index := Index{}
		index.Init()

		index.httpClient = ts.Client()

		urlSplit := strings.SplitN(ts.URL, "://", 2)

		host, portStr, _ := net.SplitHostPort(urlSplit[1])
		port, _ := strconv.Atoi(portStr)

		index.Join(net.ParseIP(host), port)
	})
}

func TestIndex_Scan(t *testing.T) {
	t.Run("Scan()", func(t *testing.T) {
		index := Index{}
		index.Init()

		index.scanner.Running = true

		index.Scan()
	})

	t.Run("Scan() => integration", func(t *testing.T) {

		if !test.Integration() {
			t.Skip()
		}

		index := Index{}
		index.Init()
		index.Scan()
	})
}

func TestIndex_Update(t *testing.T) {
	t.Run("Update()", func(t *testing.T) {
		index := Index{}
		index.Init()

		index.Update()
	})

	t.Run("Update() => Update while updating", func(t *testing.T) {
		index := Index{}
		index.Init()

		index.updateChan = make(chan bool)
		go func() {
			time.Sleep(1 * time.Nanosecond)
			close(index.updateChan)
		}()

		index.Update()
	})

	t.Run("Update() => Update while update awaiting", func(t *testing.T) {
		index := Index{}
		index.Init()

		index.updateChan = make(chan bool)
		index.updateNext = true

		index.Update()

		close(index.updateChan)
	})

	t.Run("Update() => Update with registered nodes", func(t *testing.T) {
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `{
				"Type": 0,
				"Status": 1,
				"Connections": {},
				"Nodes": {
					"ABC": {
						"ID": "ABC",
						"IP": "123.123.123.123",
						"Type": -1,
						"Port": 123
					}
				}
			}`)
		}))

		index := Index{}
		index.Init()

		index.httpClient = ts.Client()
		index.httpClient.Timeout = 100 * time.Millisecond

		urlSplit := strings.SplitN(ts.URL, "://", 2)

		host, portStr, _ := net.SplitHostPort(urlSplit[1])
		port, _ := strconv.Atoi(portStr)

		newNode := Node{
			IP:   net.ParseIP(host),
			ID:   "MOCKID",
			Type: Unknown,
			Port: port,
		}

		index.Add(&newNode)

		index.Update()
	})
}

func TestIndex_Collect(t *testing.T) {
	t.Run("Collect() => empty network", func(t *testing.T) {
		index := Index{}
		index.Init()

		network := Network{}

		index.Collect(&network)
	})

	t.Run("Collect() => expanding network", func(t *testing.T) {
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `{
				"Indexs": {
					"2523abc8780890e8849938563ab75fe8e7dafeb600dec90013b1b79b3715a2b4": {
						"Type": 0,
						"Status": 1,
						"Connections": {},
						"Nodes": {}
					}
				}
			}`)
		}))

		index := Index{}
		index.Init()

		index.httpClient = ts.Client()

		urlSplit := strings.SplitN(ts.URL, "://", 2)

		host, portStr, _ := net.SplitHostPort(urlSplit[1])
		port, _ := strconv.Atoi(portStr)

		newNode := Node{
			IP:   net.ParseIP(host),
			ID:   "MOCKID",
			Type: Unknown,
			Port: port,
		}

		index.Add(&newNode)

		network := Network{}

		index.Collect(&network)
	})

}
