package def

import (
	"testing"
)

func TestIntegration(t *testing.T) {
	t.Run("Integraction()", func(t *testing.T) {
		Integration()
	})
}

func TestFindPort(t *testing.T) {
	t.Run("FindPort()", func(t *testing.T) {
		port := GetPort()
		FindPort()
		OpenPort()

		if port != GetPort() {
			t.Errorf("Failed to revert OpenPort() lookup back to previous state.")
		}
	})
}
