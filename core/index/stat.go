package index

import (
	"runtime"
	"strconv"
	"time"

	"github.com/cloudfoundry/gosigar"
)

type Stat struct {
	Cores   int
	Memory  uint64
	Storage uint64
	Usage   StatUsage
}

type StatUsage struct {
	Cores   float64
	Memory  float64
	Storage float64
	Time    time.Time
}

type StatLoad struct {
	Short StatUsage
	Mid   StatUsage
	Long  StatUsage
	hist  []StatUsage
}

var systemLoad = StatLoad{}

func (s *Stat) getCores() int {
	return runtime.NumCPU()
}

func (s *Stat) getMemory() uint64 {
	mem := sigar.Mem{}
	mem.Get()
	return mem.Total
}

func (s *Stat) getStorage() (uint64, error) {
	_, total, err := getStorage()
	if err != nil {
		return 0, err
	}
	return total, err
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

func GetStat() (Stat, error) {
	s := Stat{}

	St, err := s.getStorage()
	if err != nil {
		return s, err
	}

	s.Cores = s.getCores()
	s.Memory = s.getMemory()
	s.Storage = St

	return s, nil
}

func GetLoad() StatLoad {
	st := time.Now()

	shortCount := 0.0
	midCount := 0.0
	longCount := 0.0

	// Reset all averages
	systemLoad.Short.Cores = 0.0
	systemLoad.Short.Memory = 0.0
	systemLoad.Short.Storage = 0.0

	systemLoad.Mid.Cores = 0.0
	systemLoad.Mid.Memory = 0.0
	systemLoad.Mid.Storage = 0.0

	systemLoad.Long.Cores = 0.0
	systemLoad.Long.Memory = 0.0
	systemLoad.Long.Storage = 0.0

	for _, cv := range systemLoad.hist {
		df := st.Sub(cv.Time)
		if df <= time.Minute {
			shortCount++
			systemLoad.Short.Cores = systemLoad.Short.Cores + ((cv.Cores - systemLoad.Short.Cores) / shortCount)
			systemLoad.Short.Memory = systemLoad.Short.Memory + ((cv.Memory - systemLoad.Short.Memory) / shortCount)
			systemLoad.Short.Storage = systemLoad.Short.Storage + ((cv.Storage - systemLoad.Short.Storage) / shortCount)
		}

		if df <= time.Hour {
			midCount++
			systemLoad.Mid.Cores = systemLoad.Mid.Cores + ((cv.Cores - systemLoad.Mid.Cores) / midCount)
			systemLoad.Mid.Memory = systemLoad.Mid.Memory + ((cv.Memory - systemLoad.Mid.Memory) / midCount)
			systemLoad.Mid.Storage = systemLoad.Mid.Storage + ((cv.Storage - systemLoad.Mid.Storage) / midCount)
		}

		if df <= time.Hour*24 {
			longCount++
			systemLoad.Long.Cores = systemLoad.Long.Cores + ((cv.Cores - systemLoad.Long.Cores) / longCount)
			systemLoad.Long.Memory = systemLoad.Long.Memory + ((cv.Memory - systemLoad.Long.Memory) / longCount)
			systemLoad.Long.Storage = systemLoad.Long.Storage + ((cv.Storage - systemLoad.Long.Storage) / longCount)
		}
	}

	return systemLoad
}

var (
	byteScale = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}
)

func bytesStr(n uint64) string {
	i := 0
	for n > 1024 {
		n = n / 1024
		i++
	}
	return strconv.FormatUint(n, 10) + byteScale[i]
}
