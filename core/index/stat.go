package index

import (
	"strconv"
	"time"
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

func (s *Stat) getStorage() (uint64, error) {
	_, total, err := getStorage()
	if err != nil {
		return 0, err
	}
	return total, err
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
