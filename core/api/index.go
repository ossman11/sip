package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/ossman11/sip/core/def"
	"github.com/ossman11/sip/core/index"
)

// Index the Api interface implementation for the Index Api
type Index struct {
	index *index.Index
}

func NewIndex() API {
	i := &index.Index{}
	i.Init()

	ret := Index{
		index: i,
	}

	return ret
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

			h.index.Add(&newNode)
		}

		if ip6 != nil {
		}

		enc.Encode(node)
	}

	if r.Method == "GET" {
		userAgent := r.Header.Get("user-agent")
		targetOS, targetArch := index.UserAgent(userAgent)

		if userAgent == "" || targetOS == "" || targetArch == "" {
			http.Error(w, "Failed to resolve user-agent platform.", http.StatusSeeOther)
			return
		}

		hostName := r.Host
		if strings.Index(r.URL.Path, "/join/"+hostName) < 0 {
			http.Redirect(w, r, r.URL.String()+"/"+hostName, http.StatusSeeOther)
			return
		}

		if targetOS == "windows" && strings.Index(r.URL.Path, "/join/"+hostName+".exe") < 0 {
			http.Redirect(w, r, r.URL.String()+".exe", http.StatusSeeOther)
			return
		}

		err := index.Build(targetOS, targetArch)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to compile target binaries.", 500)
			return
		}

		tmpFile, err := os.Open(".tmp/" + targetOS + "-" + targetArch)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to open target binaries.", 500)
			return
		}
		defer tmpFile.Close()

		io.Copy(w, tmpFile)
	}
}

func (h *Index) collect(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	network := index.Network{}

	if r.Method == "POST" {
		dec := json.NewDecoder(r.Body)
		dec.Decode(&network)
	}

	h.index.Collect(&network)

	enc.Encode(network)
}

func (h *Index) call(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	network := index.Network{}

	if r.Method == "POST" {
		dec := json.NewDecoder(r.Body)
		dec.Decode(&network)
	}

	h.index.Collect(&network)

	enc.Encode(network)
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
		def.APIIndex:             h.handleIndex,
		def.APIIndexJoin:         h.join,
		def.APIIndexJoin + "/**": h.join,
		def.APIIndexCollect:      h.collect,
		"/index/status":          h.status,
		"/index/refresh":         h.refresh,
	}
}

// Post Implements the Post API for the Index definition
func (h Index) Post() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		def.APIIndex:        h.handleIndex,
		def.APIIndexJoin:    h.join,
		def.APIIndexCollect: h.collect,
	}
}

func (h Index) Running() func() {
	return func() {
		go h.index.Scan()
	}
}
