package api

import (
	"fmt"
	"net/http"
)

const (
	// HomePageContent The text to be displayed on the home page
	HomePageContent string = "This is the home page."
)

// Home the Api interface implementation for the Home Api
type Home struct{}

func NewHome() API {
	return Home{}
}

// Action Implements the Home Api behavior
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, HomePageContent)
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
