package api

import (
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
	index *index.Index
}

func NewIndex() *Index {
	i := &index.Index{}
	i.Init()
	return &Index{
		index: i,
	}
}

// Action Implements the Index Api behavior
func (h *Index) handleIndex(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	strIP := r.RemoteAddr
	ipEnd := strings.LastIndex(strIP, ":")
	ip := net.ParseIP(strIP[:ipEnd]).To4()

	if r.Method == "POST" {
		bod := struct{ ID []string }{}
		dec := json.NewDecoder(r.Body)
		dec.Decode(&bod)

		for it := range bod.ID {
			remID := index.ParseStr(bod.ID[it])
			remID = remID.In(ip)
			bod.ID[it] = remID.String()
		}

		enc.Encode(&bod)
	} else {
		enc.Encode(h.index.GetAll(ip))
	}
}

func (h *Index) join(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	strIP := r.RemoteAddr
	ipEnd := strings.LastIndex(strIP, ":")
	ip := net.ParseIP(strIP[:ipEnd])
	ip4 := ip.To4()
	ip6 := ip.To16()

	indexType := h.index.Type

	if ip4 != nil {
		id := index.ParseIP(ip4)
		_, exists := h.index.Connections[id.String()]
		if !exists {
			h.index.Add(id.String(), index.Connection{
				IP:   ip,
				ID:   id,
				Type: index.Unknown,
			})

			res, err := http.Get("http://" + ip4.String() + ":1670" + def.APIIndexJoin)

			if err == nil {
				var m index.Type
				dec := json.NewDecoder(res.Body)
				err = dec.Decode(&m)
				if err != nil {
					fmt.Println(err)
					return
				}
				h.index.Add(id.String(), index.Connection{
					IP:   ip,
					ID:   id,
					Type: m,
				})
			}
		}
	}

	if ip6 != nil {
	}

	enc.Encode(indexType)
}

func (h *Index) status(w http.ResponseWriter, r *http.Request) {
	var s string
	if h.index.Status() {
		s = "Running"
	} else {
		s = "Idle"
	}
	fmt.Fprintf(w, s)
}

func (h *Index) refresh(w http.ResponseWriter, r *http.Request) {
	if !h.index.Status() {
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
		"/index": h.handleIndex,
	}
}

func (h Index) Running() func() {
	return func() {
		go h.index.Scan()
	}
}
