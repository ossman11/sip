// +build js

package index

import (
	"time"
)

func (s *Stat) getCores() int {
	return 1
}

func (s *Stat) getMemory() uint64 {
	return 1
}

func getStorage() (uint64, uint64, error) {
	return 0, 0, nil
}

func (s *Stat) GetUsage() (StatUsage, error) {
	u := StatUsage{}

	// Storage
	u.Storage = 1

	// Memory
	u.Memory = 1

	// Cores
	u.Cores = 1
	u.Time = time.Now()

	systemLoad.hist = append(systemLoad.hist, u)

	return u, nil
}
