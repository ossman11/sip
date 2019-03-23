package index

import (
	"fmt"
	"strings"
	"testing"
)

func TestExtendRoute(t *testing.T) {
	t.Run("ExtendRoute()", func(t *testing.T) {
		r := NewRoute([]ID{NewIDHash("Node1"), NewIDHash("Node2")})
		e := ExtendRoute(NewIDHash("Node3"), &r)

		if e.Equal(&r) {
			t.Error("The extended Route was equal to the original Route.")
		}
	})
}

func TestRoute_Next(t *testing.T) {
	t.Run("Next()", func(t *testing.T) {
		route := NewRoute([]ID{
			NewID("1"),
			NewID("2"),
			NewID("3"),
			NewID("4"),
		})

		nextRoute := route.Next()

		fmt.Println(nextRoute.String())
	})
}

func TestRoute_String(t *testing.T) {
	t.Run("Next()", func(t *testing.T) {
		route := NewRoute([]ID{
			NewID("1"),
			NewID("2"),
			NewID("3"),
			NewID("4"),
		})
		str := route.String()

		strPath := strings.Split(str, ",")
		path := NewRoute([]ID{})
		for _, pv := range strPath {
			path = NewRoute(append(path.Nodes, NewID(pv)))
		}

		fmt.Println(str)
	})
}
