// +build !js

package index

import (
	"runtime"
	"time"

	"github.com/cloudfoundry/gosigar"
)

func (s *Stat) getCores() int {
	return runtime.NumCPU()
}

func (s *Stat) getMemory() uint64 {
	mem := sigar.Mem{}
	mem.Get()
	return mem.Total
}

func (s *Stat) GetUsage() (StatUsage, error) {
	u := StatUsage{}

	// Storage
	free, total, err := getStorage()
	if err != nil {
		return u, err
	}
	u.Storage = float64(total-free) / float64(total)

	// Memory
	mem := sigar.Mem{}
	mem.Get()
	u.Memory = float64(mem.Used) / float64(mem.Total)

	// Cores
	cpu1 := sigar.Cpu{}
	cpu1.Get()
	time1 := time.Now()

	time.Sleep(100 * time.Millisecond)

	cpu2 := sigar.Cpu{}
	cpu2.Get()
	time2 := time.Now()

	del := cpu2.Delta(cpu1)
	cpuUsage := float64(del.Total()-del.Idle) / float64(time2.Sub(time1)) / float64(runtime.NumCPU())
	u.Cores = cpuUsage

	u.Time = time1.Add(time2.Sub(time1) / 2)

	systemLoad.hist = append(systemLoad.hist, u)

	return u, nil
}
