package core

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ossman11/sip/core/api"
)

// Server interface for the core server functionalities
type Server interface {
	Start()
}

// Core implements the root functionality of a sip instance
type Core struct {
	ready   bool
	busy    *sync.Mutex
	server  *http.Server
	handler *api.Handler
}

// Init ensures that the Core instance is in a prepared state
func (c *Core) Init() {
	if c.busy == nil {
		c.busy = &sync.Mutex{}
	}
	c.busy.Lock()
	// Prevent overwriting the existing instance
	if c.ready {
		c.busy.Unlock()
		return
	}

	// Create new instance for the handler and server
	c.handler = &api.Handler{}
	c.handler.AddCore()
	c.server = &http.Server{
		Addr:           ":1670",
		Handler:        c.handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	c.handler.Running()
	c.ready = true
	c.busy.Unlock()
}

// Starts the core instance
func (c *Core) Start() {
	c.Init()
	log.Fatal(c.server.ListenAndServe())
}

// AddApis enhances the Core functionalities with additional Api implementations
func (c Core) AddApis(n api.API) {
	c.Init()
	c.handler.Add(n)
}

// NewServer creates a new Core instance Server interface
func NewServer() *Core {
	r := &Core{}
	r.Init()
	return r
}
