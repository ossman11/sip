package api

import "net/http"

// API provides the interface through which it is possible to expose API on the Core Server
type API interface {
	Get() map[string]http.HandlerFunc
	Post() map[string]http.HandlerFunc
	Running() func()
}
