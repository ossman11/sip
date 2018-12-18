package test

import (
	"testing"

	"github.com/ossman11/sip/core/def"
)

func TestIntegration(t *testing.T) {
	t.Run("Integraction()", func(t *testing.T) {
		Integration()
	})
}

func TestFindPort(t *testing.T) {
	t.Run("FindPort()", func(t *testing.T) {
		port := def.GetPort()
		FindPort()
		OpenPort()

		if port != def.GetPort() {
			t.Errorf("Failed to revert OpenPort() lookup back to previous state.")
		}
	})
}
