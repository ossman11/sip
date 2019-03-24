// +build js

package main

import (
	"bytes"
	"encoding/json"
	"net"
	"strconv"
	"syscall/js"
	"time"

	"github.com/ossman11/sip/core/def"

	"github.com/ossman11/sip/core/index"
)

func main() {
	global := js.Global()
	globSip := ensure(global, "sip")

	i := index.Index{}
	i.Init()

	location := global.Get("location")
	hostname := location.Get("hostname")
	port := location.Get("port")
	portInt, err := strconv.Atoi(port.String())

	if err != nil {
		portInt = def.GetPort()
	}

	hostIP := net.ParseIP(hostname.String())

	if hostIP == nil {
		hostIP = net.ParseIP("127.0.0.1")
	}

	i.Join(hostIP, portInt)

	iJS, err := object(i)

	if err != nil {
		panic(err)
	}

	globSip.Set("index", iJS)

	callFnc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			t := args[0]
			if t.Type() != js.TypeString {
				return // js.ValueOf("Failed not enough arguments.")
			}

			n := index.Network{}
			i.Collect(&n)
			s := index.ThisNode(&i, hostIP)
			err, ps := n.Path(s.ID, index.NewID(t.String()))

			if err != nil {
				return // js.ValueOf("Failed no path found.")
			}

			// Picks random route from all options
			p := index.Route{}
			for _, cv := range ps {
				p = *cv[0]
				continue
			}

			i.Call(p, "index")
		}()
		return nil
	})

	globSip.Set("call", callFnc)

	collectFnc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			// Collect network
			n := index.Network{}
			i.Collect(&n)

			// Tranform to string
			var b bytes.Buffer
			bod, _ := json.Marshal(n)
			b.Write(bod)

			// Convert to javascript JSON
			jsJSON := js.Global().Get("JSON").Call("parse", b.String())
			js.Global().Get("console").Call("log", jsJSON)
		}()
		return nil
	})

	globSip.Set("collect", collectFnc)

	// Keep Go alive in the background
	go forever()
	select {}
}

// Ensures that the property exists
func ensure(obj js.Value, prop string) js.Value {
	org := obj.Get(prop)
	if org.Type() == js.TypeUndefined {
		obj.Set(prop, map[string]interface{}{})
		return obj.Get(prop)
	}
	return org
}

func object(v interface{}) (interface{}, error) {
	var ret interface{}

	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func forever() {
	for {
		time.Sleep(time.Hour * 24)
	}
}
