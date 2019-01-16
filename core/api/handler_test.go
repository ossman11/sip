package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	getSuccess      string = "GET SUCCESS"
	postSuccess     string = "POST SUCCESS"
	notFoundSuccess string = "404 page not found\n"
)

type GetHandler struct{}

func (h *GetHandler) GetHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, getSuccess)
}

func (h GetHandler) Get() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/GET": h.GetHandlerFunc,
	}
}

func (h GetHandler) Post() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{}
}

func (h GetHandler) Running() func() {
	return func() {}
}

func NewGetHandler() API {
	ret := GetHandler{}
	return ret
}

type PostHandler struct{}

func (h *PostHandler) PostHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, postSuccess)
}

func (h PostHandler) Get() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{}
}

func (h PostHandler) Post() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/POST": h.PostHandlerFunc,
	}
}

func (h PostHandler) Running() func() {
	return func() {}
}

func NewPostHandler() API {
	ret := PostHandler{}
	return ret
}

func TestHandler_AddCore(t *testing.T) {

	t.Run("AddCore() => default", func(t *testing.T) {
		handler := Handler{}
		handler.AddCore()
	})

	t.Run("AddCore() => default running", func(t *testing.T) {
		handler := Handler{}
		handler.AddCore()
		handler.Running()
	})

	t.Run("AddCore() => missing running", func(t *testing.T) {
		handler := Handler{}
		handler.Running()
	})

	t.Run("AddCore() => duplicate GET", func(t *testing.T) {
		orgAPIs := coreAPIs
		coreAPIs = []func() API{
			NewGetHandler,
			NewGetHandler,
		}

		handler := Handler{}
		handler.AddCore()

		coreAPIs = orgAPIs
	})

	t.Run("AddCore() => duplicate POST", func(t *testing.T) {
		orgAPIs := coreAPIs
		coreAPIs = []func() API{
			NewPostHandler,
			NewPostHandler,
		}

		handler := Handler{}
		handler.AddCore()

		coreAPIs = orgAPIs
	})
}

func TestHandler_ServeHTTP(t *testing.T) {

	t.Run("ServeHTTP() => GET request", func(t *testing.T) {
		orgAPIs := coreAPIs
		coreAPIs = []func() API{
			NewGetHandler,
		}

		handler := Handler{}
		handler.AddCore()

		coreAPIs = orgAPIs

		req := httptest.NewRequest("GET", "http://localhost/GET", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			t.Errorf("Failed to use ServeHTTP(), because StatusCode was %v.", resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Errorf("Failed to use ServeHTTP(), because Content-Type was %v.", resp.Header.Get("Content-Type"))
		}

		if string(body) != getSuccess {
			t.Errorf("Failed to use ServeHTTP(), because Body was %v.", string(body))
		}
	})

	t.Run("ServeHTTP() => POST request", func(t *testing.T) {
		orgAPIs := coreAPIs
		coreAPIs = []func() API{
			NewPostHandler,
		}

		handler := Handler{}
		handler.AddCore()

		coreAPIs = orgAPIs

		req := httptest.NewRequest("POST", "http://localhost/POST", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			t.Errorf("Failed to use ServeHTTP(), because StatusCode was %v.", resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Errorf("Failed to use ServeHTTP(), because Content-Type was %v.", resp.Header.Get("Content-Type"))
		}

		if string(body) != postSuccess {
			t.Errorf("Failed to use ServeHTTP(), because Body was %v.", string(body))
		}
	})

	t.Run("ServeHTTP() => Not found GET request", func(t *testing.T) {
		orgAPIs := coreAPIs
		coreAPIs = []func() API{}

		handler := Handler{}
		handler.AddCore()

		coreAPIs = orgAPIs

		req := httptest.NewRequest("GET", "http://localhost/NOTFOUND", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != 404 {
			t.Errorf("Failed to use ServeHTTP(), because StatusCode was %v.", resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Errorf("Failed to use ServeHTTP(), because Content-Type was %v.", resp.Header.Get("Content-Type"))
		}

		if string(body) != notFoundSuccess {
			t.Errorf("Failed to use ServeHTTP(), because Body was %v.", string(body))
		}
	})

	t.Run("ServeHTTP() => Not found POST request", func(t *testing.T) {
		orgAPIs := coreAPIs
		coreAPIs = []func() API{}

		handler := Handler{}
		handler.AddCore()

		coreAPIs = orgAPIs

		req := httptest.NewRequest("POST", "http://localhost/NOTFOUND", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != 404 {
			t.Errorf("Failed to use ServeHTTP(), because StatusCode was %v.", resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Errorf("Failed to use ServeHTTP(), because Content-Type was %v.", resp.Header.Get("Content-Type"))
		}

		if string(body) != notFoundSuccess {
			t.Errorf("Failed to use ServeHTTP(), because Body was %v.", string(body))
		}
	})

	t.Run("ServeHTTP() => GET core Home page", func(t *testing.T) {
		handler := Handler{}
		handler.AddCore()

		req := httptest.NewRequest("GET", "http://localhost/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			t.Errorf("Failed to use ServeHTTP(), because StatusCode was %v.", resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Errorf("Failed to use ServeHTTP(), because Content-Type was %v.", resp.Header.Get("Content-Type"))
		}

		if string(body) != HomePageContent {
			t.Errorf("Failed to use ServeHTTP(), because Body was %v.", string(body))
		}
	})
}

func TestHandler_Add(t *testing.T) {
	t.Run("Add() => Empty API", func(t *testing.T) {
		handler := Handler{}
		handler.Add(Empty{})
	})

	t.Run("Add() => Empty API running", func(t *testing.T) {
		handler := Handler{}
		handler.Add(Empty{})
		handler.Running()
	})
}
