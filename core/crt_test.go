package core

import (
	"os"
	"testing"
)

func TestGenCrt(t *testing.T) {
	t.Run("GenCrt()", func(t *testing.T) {
		// Remove certificate if they exist
		os.Remove("./crt/server.crt")
		os.Remove("./crt/server.key")

		// Generate certificate
		err := GenCrt()
		if err != nil {
			t.Error(err)
		}

		// Ensure caching is working
		err = GenCrt()
		if err != nil {
			t.Error(err)
		}
	})
}
