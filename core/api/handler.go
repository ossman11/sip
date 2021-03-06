package api

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
)

// Handler implements the http handeling of a sip instance
type Handler struct {
	gets     map[string]http.HandlerFunc
	posts    map[string]http.HandlerFunc
	Runnings []func()
}

func (rh *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l := rh.gets
	if r.Method == "POST" {
		l = rh.posts
	}

	a, ok := l[r.URL.Path]

	// Set default headers for all services (e.g. security)
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Header().Add("Access-Control-Allow-Origin", r.Host)
	w.Header().Add("Access-Control-Allow-Credentials", "false")
	w.Header().Add("Access-Control-Allow-Headers", "Client-Protocol, Content-Length, Content-Type, x-target")

	if !ok {
		// search through glob matching
		for ck, cv := range l {
			if m, err := filepath.Match(ck, r.URL.Path); m && err == nil {
				a = cv
				ok = true
				break
			}
		}
	}

	if ok {
		a(w, r)
		return
	}
	http.NotFound(w, r)

}

// Add adds a map of Api implementations to the Handler instance
func (rh *Handler) Add(a API) error {
	// Ensure that the gets map is initialized
	if rh.gets == nil {
		rh.gets = map[string]http.HandlerFunc{}
	}
	// Ensure that the posts map is initialized
	if rh.posts == nil {
		rh.posts = map[string]http.HandlerFunc{}
	}

	// Add all the get urls of the current API defention
	g := a.Get()
	for k := range g {
		_, ok := rh.gets[k]
		if ok {
			return errors.New("Failed to add api: \"" + k + "\", because this api is already assigned.")
		}
		rh.gets[k] = g[k]
	}

	// Add all the post urls of the current API defention
	p := a.Post()
	for k := range p {
		_, ok := rh.posts[k]
		if ok {
			return errors.New("Failed to add api: \"" + k + "\", because this api is already assigned.")
		}
		rh.posts[k] = p[k]
	}

	rh.Runnings = append(rh.Runnings, a.Running())

	return nil
}

func (rh *Handler) Running() {
	if rh.Runnings == nil {
		return
	}

	for _, r := range rh.Runnings {
		r()
	}
}

var coreAPIs = []func() API{
	NewIndex,
	NewHome,
}

// AddCore adds the core API implementations to the Handler instance
func (rh *Handler) AddCore() {
	for _, c := range coreAPIs {
		err := rh.Add(c())
		if err != nil {
			fmt.Println(err)
		}
	}
}
