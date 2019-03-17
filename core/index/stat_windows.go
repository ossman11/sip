// +build windows

package index

import (
	"os"
	"syscall"
	"unsafe"
)

func getStorage() (uint64, uint64, error) {
	wd, err := os.Getwd()
	if err != nil {
		return 0, 0, err
	}

	h := syscall.MustLoadDLL("kernel32.dll")
	c := h.MustFindProc("GetDiskFreeSpaceExW")

	var userBytes uint64
	var totalBytes uint64
	var freeBytes uint64

	_, _, err = c.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(wd))),
		uintptr(unsafe.Pointer(&userBytes)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&freeBytes)),
	)
	if err != nil && err.Error() != "The operation completed successfully." {
		return 0, 0, err
	}

	return userBytes, totalBytes, nil
}
