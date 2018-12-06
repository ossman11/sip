package index

import (
	"testing"
	"time"
)

func TestScan_awaitChan(t *testing.T) {
	t.Run("awaitChan()", func(t *testing.T) {
		scan := Scan{}
		cc, ch := scan.startChan(0)

		go func() {
			time.Sleep(1 * time.Nanosecond)
			scan.endGo(ch)
		}()
		cc++

		scan.awaitChan(cc, ch)
	})
}
