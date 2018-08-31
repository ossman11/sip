package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/ossman11/sip/core/def"
	"github.com/ossman11/sip/core/index"
)

// Index the Api interface implementation for the Index Api
type Index struct {
	index      *index.Index
	httpClient *http.Client
}

func NewIndex() *Index {
	i := &index.Index{}
	i.Init()

	// Always scan without security enabled
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &Index{
		index: i,
		httpClient: &http.Client{
			Transport: tr,
		},
	}
}

// Action Implements the Index Api behavior
func (h *Index) handleIndex(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	/*
		strIP := r.RemoteAddr
		ipEnd := strings.LastIndex(strIP, ":")
		ip := net.ParseIP(strIP[:ipEnd]).To4()
	*/
	if r.Method == "POST" {
		bod := index.Index{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&bod)

		if err != nil {
			fmt.Println(err)
		}

		if h.index.Merge(&bod) {
			go h.index.Update()
		}

		enc.Encode(h.index)
	} else {
		enc.Encode(h.index)
	}
}

func (h *Index) join(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	strIP := r.RemoteAddr
	ipEnd := strings.LastIndex(strIP, ":")
	ip := net.ParseIP(strIP[:ipEnd])
	// fmt.Println("Join from ip: ", ip)
	ip4 := ip.To4()
	ip6 := ip.To16()

	node := index.ThisNode(h.index, ip)

	if r.Method == "POST" {
		if ip4 != nil {
			newNode := index.Node{}
			dec := json.NewDecoder(r.Body)
			dec.Decode(&newNode)
			// newNode.IP = ip4

			h.index.JoinNode(newNode)
		}

		if ip6 != nil {
		}
	}

	enc.Encode(node)
}

func (h *Index) status(w http.ResponseWriter, r *http.Request) {
	s := h.index.Status.String()
	fmt.Fprintf(w, s)
}

func (h *Index) refresh(w http.ResponseWriter, r *http.Request) {
	if h.index.Status == index.Idle {
		go h.index.Scan()
	}
	fmt.Fprintf(w, "")
}

// Get Implements the Get API for the Index definition
func (h Index) Get() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		def.APIIndex:     h.handleIndex,
		def.APIIndexJoin: h.join,
		"/index/status":  h.status,
		"/index/refresh": h.refresh,
	}
}

// Post Implements the Post API for the Index definition
func (h Index) Post() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		def.APIIndex:     h.handleIndex,
		def.APIIndexJoin: h.join,
	}
}

func (h Index) Running() func() {
	return func() {
		go h.index.Scan()
	}
}
