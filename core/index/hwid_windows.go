// +build windows

package index

import (
	"golang.org/x/sys/windows/registry"
)

// Credit: https://github.com/denisbrodbeck/machineid/blob/master/id_windows.go

func fetchHWID() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE|registry.WOW64_64KEY)
	if err != nil {
		return ""
	}
	defer k.Close()

	s, _, err := k.GetStringValue("MachineGuid")
	if err != nil {
		return ""
	}
	return s
}
