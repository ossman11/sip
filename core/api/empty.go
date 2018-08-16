package api

import (
	"net/http"
)

// Empty the Api interface implementation for the Empty Api
type Empty struct{}

// Get Implements the Get API for the Empty definition
func (h Empty) Get() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{}
}

// Post Implements the Post API for the Empty definition
func (h Empty) Post() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{}
}
