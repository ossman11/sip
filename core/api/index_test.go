package api

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ossman11/sip/core/index"
)

func NewIndexTest() Index {
	i := &index.Index{}
	i.Init()

	ret := Index{
		index: i,
	}

	return ret
}
func TestIndex_handleIndex(t *testing.T) {

	t.Run("handleIndex() => GET", func(t *testing.T) {
		testIndex := NewIndexTest()

		req := httptest.NewRequest("GET", "http://localhost/", nil)
		w := httptest.NewRecorder()
		testIndex.handleIndex(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			t.Errorf("Failed to use handleIndex(), because StatusCode was %v.", resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Errorf("Failed to use handleIndex(), because Content-Type was %v.", resp.Header.Get("Content-Type"))
		}

		e := httptest.NewRecorder()
		enc := json.NewEncoder(e)
		enc.Encode(testIndex.index)

		expt := e.Result()
		exptBytes, _ := ioutil.ReadAll(expt.Body)

		if string(body) != string(exptBytes) {
			t.Errorf("Failed to use handleIndex(), because Body was %v.", string(body))
		}
	})

	t.Run("handleIndex() => POST(empty index)", func(t *testing.T) {
		testIndex := NewIndexTest()

		req := httptest.NewRequest("POST", "http://localhost/", strings.NewReader(`{
			"Type": 0,
			"Status": 1,
			"Connections": {},
			"Nodes": {}
		}`))
		w := httptest.NewRecorder()
		testIndex.handleIndex(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			t.Errorf("Failed to use handleIndex(), because StatusCode was %v.", resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Errorf("Failed to use handleIndex(), because Content-Type was %v.", resp.Header.Get("Content-Type"))
		}

		e := httptest.NewRecorder()
		enc := json.NewEncoder(e)
		enc.Encode(testIndex.index)

		expt := e.Result()
		exptBytes, _ := ioutil.ReadAll(expt.Body)

		if string(body) != string(exptBytes) {
			t.Errorf("Failed to use handleIndex(), because Body was %v.", string(body))
		}
	})

	t.Run("handleIndex() => POST(additional index)", func(t *testing.T) {
		testIndex := NewIndexTest()

		req := httptest.NewRequest("POST", "http://localhost/", strings.NewReader(`{
			"Type": 0,
			"Status": 1,
			"Connections": {},
			"Nodes": {
				"abc": {
					"IP": "127.0.0.1",
					"ID": "abc",
					"Type": 0,
					"Port": 1670
				}
			}
		}`))
		w := httptest.NewRecorder()
		testIndex.handleIndex(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			t.Errorf("Failed to use handleIndex(), because StatusCode was %v.", resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Errorf("Failed to use handleIndex(), because Content-Type was %v.", resp.Header.Get("Content-Type"))
		}

		e := httptest.NewRecorder()
		enc := json.NewEncoder(e)
		enc.Encode(testIndex.index)

		expt := e.Result()
		exptBytes, _ := ioutil.ReadAll(expt.Body)

		if string(body) != string(exptBytes) {
			t.Errorf("Failed to use handleIndex(), because Body was %v.", string(body))
		}
	})

	t.Run("handleIndex() => POST(invalid index)", func(t *testing.T) {
		testIndex := NewIndexTest()

		req := httptest.NewRequest("POST", "http://localhost/", strings.NewReader(`{Nodes:[]}`))
		w := httptest.NewRecorder()
		testIndex.handleIndex(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			t.Errorf("Failed to use handleIndex(), because StatusCode was %v.", resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Errorf("Failed to use handleIndex(), because Content-Type was %v.", resp.Header.Get("Content-Type"))
		}

		e := httptest.NewRecorder()
		enc := json.NewEncoder(e)
		enc.Encode(testIndex.index)

		expt := e.Result()
		exptBytes, _ := ioutil.ReadAll(expt.Body)

		if string(body) != string(exptBytes) {
			t.Errorf("Failed to use handleIndex(), because Body was %v.", string(body))
		}
	})

}

func TestIndex_join(t *testing.T) {
	t.Run("join() => GET", func(t *testing.T) {
		testIndex := NewIndexTest()

		req := httptest.NewRequest("GET", "http://localhost/", nil)
		w := httptest.NewRecorder()
		testIndex.join(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Failed to use join(), because StatusCode was %v.", resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Errorf("Failed to use join(), because Content-Type was %v.", resp.Header.Get("Content-Type"))
		}

		expt := "Failed to resolve user-agent platform."
		bodyStr := strings.Trim(string(body), " \n")

		if bodyStr != expt {
			t.Errorf("Failed to use join(), because Body was:\n  %v\n  Expecting:\n  %v.", bodyStr, expt)
		}
	})

	t.Run("join() => POST", func(t *testing.T) {
		testIndex := NewIndexTest()

		req := httptest.NewRequest("POST", "http://localhost/", strings.NewReader(`{}`))
		w := httptest.NewRecorder()
		testIndex.join(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			t.Errorf("Failed to use join(), because StatusCode was %v.", resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Errorf("Failed to use join(), because Content-Type was %v.", resp.Header.Get("Content-Type"))
		}

		e := httptest.NewRecorder()
		enc := json.NewEncoder(e)

		strIP := req.RemoteAddr
		ipEnd := strings.LastIndex(strIP, ":")
		ip := net.ParseIP(strIP[:ipEnd])
		ip4 := ip.To4()

		node := index.ThisNode(testIndex.index, ip4)
		enc.Encode(node)

		expt := e.Result()
		exptBytes, _ := ioutil.ReadAll(expt.Body)

		if string(body) != string(exptBytes) {
			t.Errorf("Failed to use join(), because Body was:\n  %v\n  Expecting:\n  %v.", string(body), string(exptBytes))
		}
	})

}

func TestIndex_refresh(t *testing.T) {
	t.Run("refresh() => GET", func(t *testing.T) {
		testIndex := NewIndexTest()

		req := httptest.NewRequest("GET", "http://localhost/", nil)
		w := httptest.NewRecorder()
		testIndex.refresh(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			t.Errorf("Failed to use refresh(), because StatusCode was %v.", resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Errorf("Failed to use refresh(), because Content-Type was %v.", resp.Header.Get("Content-Type"))
		}

		if string(body) != "" {
			t.Errorf("Failed to use refresh(), because Body was:\n  %v\n  Expecting:\n  %v.", string(body), "")
		}
	})
}

func TestIndex_status(t *testing.T) {
	t.Run("status() => GET", func(t *testing.T) {
		testIndex := NewIndexTest()

		req := httptest.NewRequest("GET", "http://localhost/", nil)
		w := httptest.NewRecorder()
		testIndex.status(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			t.Errorf("Failed to use status(), because StatusCode was %v.", resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Errorf("Failed to use status(), because Content-Type was %v.", resp.Header.Get("Content-Type"))
		}

		if string(body) != testIndex.index.Status.String() {
			t.Errorf("Failed to use status(), because Body was:\n  %v\n  Expecting:\n  %v.", string(body), testIndex.index.Status.String())
		}
	})
}

func TestIndex_collect(t *testing.T) {
	t.Run("collect() => GET", func(t *testing.T) {
		testIndex := NewIndexTest()

		req := httptest.NewRequest("GET", "http://localhost/", nil)
		w := httptest.NewRecorder()
		testIndex.collect(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			t.Errorf("Failed to use collect(), because StatusCode was %v.", resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Errorf("Failed to use collect(), because Content-Type was %v.", resp.Header.Get("Content-Type"))
		}

		network := index.Network{}
		testIndex.index.Collect(&network)

		e := httptest.NewRecorder()
		enc := json.NewEncoder(e)
		enc.Encode(network)

		expt := e.Result()
		exptBytes, _ := ioutil.ReadAll(expt.Body)

		if string(body) != string(exptBytes) {
			t.Errorf("Failed to use collect(), because Body was:\n  %v\n  Expecting:\n  %v.", string(body), "")
		}
	})

	t.Run("collect() => POST", func(t *testing.T) {
		testIndex := NewIndexTest()

		req := httptest.NewRequest("POST", "http://localhost/", strings.NewReader(`{}`))
		w := httptest.NewRecorder()
		testIndex.collect(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			t.Errorf("Failed to use collect(), because StatusCode was %v.", resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Errorf("Failed to use collect(), because Content-Type was %v.", resp.Header.Get("Content-Type"))
		}

		network := index.Network{}
		testIndex.index.Collect(&network)

		e := httptest.NewRecorder()
		enc := json.NewEncoder(e)
		enc.Encode(network)

		expt := e.Result()
		exptBytes, _ := ioutil.ReadAll(expt.Body)

		if string(body) != string(exptBytes) {
			t.Errorf("Failed to use collect(), because Body was:\n  %v\n  Expecting:\n  %v.", string(body), "")
		}
	})
}
