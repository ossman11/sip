package core

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ossman11/sip/core/api"
	"github.com/ossman11/sip/core/def"
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

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			// tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	// Create new instance for the handler and server
	c.handler = &api.Handler{}
	c.handler.AddCore()
	fmt.Println(def.GetPort())
	c.server = &http.Server{
		Addr:           ":" + strconv.Itoa(def.GetPort()),
		Handler:        c.handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      cfg,
	}
	c.ready = true
	c.busy.Unlock()
}

// Starts the core instance
func (c *Core) Start() {
	c.Init()
	c.handler.Running()
	log.Fatal(c.server.ListenAndServeTLS("crt/server.crt", "crt/server.key"))
	// log.Fatal(c.server.ListenAndServe())
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
