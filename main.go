package main

import (
	"github.com/ossman11/sip/core"
)

func main() {
	c := core.NewServer()
	c.Start()
}
