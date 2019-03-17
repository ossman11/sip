package index

import (
	"fmt"
	"testing"
	"time"
)

func TestGetStats(t *testing.T) {
	t.Run("GetStats()", func(t *testing.T) {
		s, err := GetStat()
		if err != nil {
			t.Error(err)
		}
		lt := time.Now()
		for {
			_, err := s.GetUsage()
			if err != nil {
				t.Error(err)
			}
			if time.Now().Sub(lt) > time.Minute {
				lt = time.Now()
				l := GetLoad()
				fmt.Printf("Short Cores: %.2f Memory: %.2f Storage: %.2f", l.Short.Cores, l.Short.Memory, l.Short.Storage)
				fmt.Println()
				fmt.Printf("Mid   Cores: %.2f Memory: %.2f Storage: %.2f", l.Mid.Cores, l.Mid.Memory, l.Mid.Storage)
				fmt.Println()
				fmt.Printf("Long  Cores: %.2f Memory: %.2f Storage: %.2f", l.Long.Cores, l.Long.Memory, l.Long.Storage)
				fmt.Println()
				fmt.Println("===============================================")
			}
			/*
				fmt.Printf("Cores: %.2f Memory: %.2f Storage: %.2f", u.Cores, u.Memory, u.Storage)
				fmt.Println()
			*/
		}
	})
}
