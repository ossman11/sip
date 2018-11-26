package def

import (
	"os"
	"strconv"
	"testing"
)

func TestGetPort(t *testing.T) {

	t.Run("GetPort() => default", func(t *testing.T) {
		// Clear cache
		curPort = 0

		res := GetPort()
		if res != Port {
			t.Errorf("GetPort() returned: %v, expected %v", res, Port)
		}
	})

	t.Run("GetPort() => custom", func(t *testing.T) {
		// Clear cache
		curPort = 0

		// Set custom enviroment variable
		exp := 987654
		os.Setenv("PORT", strconv.Itoa(exp))
		res := GetPort()
		if res != exp {
			t.Errorf("GetPort() returned: %v, expected %v", res, exp)
		}
	})

	t.Run("GetPort() => custom string", func(t *testing.T) {
		// Clear cache
		curPort = 0

		// Set custom enviroment variable
		os.Setenv("PORT", "NOT A NUMBER")

		res := GetPort()
		if res != Port {
			t.Errorf("GetPort() returned: %v, expected %v", res, Port)
		}
	})
}
