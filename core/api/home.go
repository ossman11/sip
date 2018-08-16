package api

import (
	"fmt"
	"net/http"
)

// Home the Api interface implementation for the Home Api
type Home struct{}

func NewHome() *Home {
	return &Home{}
}

// Action Implements the Home Api behavior
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the home page.")
}

// Get Implements the Get API for the Home definition
func (h Home) Get() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/": home,
	}
}

// Post Implements the Post API for the Home definition
func (h Home) Post() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{}
}

func (h Home) Running() func() {
	return func() {}
}
