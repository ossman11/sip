package index

import (
	"testing"
)

func TestNewID(t *testing.T) {
	t.Run("NewID()", func(t *testing.T) {
		src := "TestID"
		id := NewID(src)

		if string(id) != src {
			t.Error("Failed to maintain original text used to create ID.")
		}
	})
}

func TestNewIDHash(t *testing.T) {
	t.Run("NewIDHash()", func(t *testing.T) {
		src := "TestID"
		id := NewIDHash(src)

		if string(id) == src {
			t.Error("Failed to hash original text used to create ID.")
		}
	})
}

func TestNewIDRoute(t *testing.T) {
	t.Run("NewIDRoute() => Correct", func(t *testing.T) {
		src := []ID{NewIDHash("1"), NewIDHash("2"), NewIDHash("3")}
		id := NewIDRoute(src)

		if id == src[0] {
			t.Error("Failed to hash original text used to create ID.")
		}
	})

	t.Run("NewIDRoute() => Empty", func(t *testing.T) {
		src := []ID{}
		id := NewIDRoute(src)

		if id != NewIDHash("") {
			t.Error("Failed to fallback to default empty Hash ID.")
		}
	})
}
