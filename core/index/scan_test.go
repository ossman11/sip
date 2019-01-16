package index

import (
	"errors"
	"net"
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

func TestScan_Scan(t *testing.T) {
	t.Run("Scan() => Fail to fetch interfaces", func(t *testing.T) {
		index := Index{}
		index.Init()

		scan := NewScan(&index)
		scan.Scan()
	})

	t.Run("Scan() => Fail to fetch interfaces", func(t *testing.T) {
		getInterfaces = func() ([]net.Interface, error) {
			return nil, errors.New("Mock error")
		}

		index := Index{}
		index.Init()

		scan := NewScan(&index)
		scan.Scan()

		getInterfaces = net.Interfaces
	})
}

func TestScan_walkIP(t *testing.T) {
	type args struct {
		ipnet *net.IPNet
		c     chan bool
	}
	tests := []struct {
		name string
		i    *Scan
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.i.walkIP(tt.args.ipnet, tt.args.c)
		})
	}
}
