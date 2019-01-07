package index

import (
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
