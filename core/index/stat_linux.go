// +build linux darwin freebsd netbsd openbsd dragonfly

package index

import (
	"os"
	"syscall"
)

func getStorage() (uint64, uint64, error) {
	var stat syscall.Statfs_t
	wd, err := os.Getwd()
	if err != nil {
		return 0, 0, err
	}
	syscall.Statfs(wd, &stat)
	return stat.Bavail * uint64(stat.Bsize), stat.Blocks * uint64(stat.Bsize), nil
}
