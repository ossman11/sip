package def

import (
	"testing"
)

func checkServer() bool {
	httpClient := HttpClient()
	httpServer := HttpServer()

	_, err := httpClient.Get(httpServer)

	if err == nil {
		return true
	}
	return false
}

func TestIntegration(t *testing.T) {
	t.Run("Integraction() => external", func(t *testing.T) {
		ready = false

		r := Integration()

		if *integration {
			if r {
				return
			}

			if !checkServer() {
				return
			}

			t.Error("Failed to determine correct integration state.", r)
		}

		if r {
			t.Error("Failed to determine correct integration state.", r)
		}
	})

	t.Run("Integraction() => forced false", func(t *testing.T) {
		ready = false
		orgIntegration := integration
		defer func() { integration = orgIntegration }()

		b := false
		integration = &b

		r := Integration()

		if r {
			t.Error("Failed to determine correct integration state.", r)
		}
	})

	t.Run("Integraction() => forced true", func(t *testing.T) {
		ready = false
		orgIntegration := integration
		defer func() { integration = orgIntegration }()

		b := true
		integration = &b

		r := Integration()

		if r && !checkServer() {
			t.Error("Failed to determine correct integration state.", r)
		}
	})

	t.Run("Integraction() => forced ready", func(t *testing.T) {
		ready = true
		orgIntegration := integration
		defer func() { integration = orgIntegration }()

		b := true
		integration = &b

		r := Integration()

		if !r {
			t.Error("Failed to determine correct integration state.", r)
		}
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
