// +build windows

package index

import (
	"syscall"
	"unsafe"
)

// Credit: https://github.com/denisbrodbeck/machineid/blob/master/id_windows.go
// Credit: https://github.com/golang/sys/blob/master/windows/registry/key.go#L78

func fetchHWID() string {
	var access uint32
	k := syscall.HKEY_LOCAL_MACHINE
	path := `SOFTWARE\Microsoft\Cryptography`
	access = 0x00001 | 0x00100

	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return ""
	}
	var subkey syscall.Handle
	err = syscall.RegOpenKeyEx(syscall.Handle(k), p, 0, access, &subkey)
	if err != nil {
		return ""
	}
	defer func() {
		syscall.RegCloseKey(subkey)
	}()

	// getValue
	buf := make([]byte, 64)
	p, err = syscall.UTF16PtrFromString("MachineGuid")
	if err != nil {
		return ""
	}
	var t uint32
	n := uint32(len(buf))
	for {
		err = syscall.RegQueryValueEx(subkey, p, nil, &t, (*byte)(unsafe.Pointer(&buf[0])), &n)
		if err == nil {
			break
		}
		if err != syscall.ERROR_MORE_DATA {
			return ""
		}
		if n <= uint32(len(buf)) {
			return ""
		}
		buf = make([]byte, n)
	}
	u := (*[1 << 29]uint16)(unsafe.Pointer(&buf[0]))[:]
	return syscall.UTF16ToString(u)
}
