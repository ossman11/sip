// +build freebsd netbsd openbsd dragonfly

package index

import (
	"bytes"
	"os"
	"strings"
)

// Credit: https://github.com/denisbrodbeck/machineid/blob/master/id_bsd.go

const hostidPath = "/etc/hostid"

// machineID returns the uuid specified at `/etc/hostid`.
// If the returned value is empty, the uuid from a call to `kenv -q smbios.system.uuid` is returned.
// If there is an error an empty string is returned.
func fetchHWID() string {
	id, err := readHostid()
	if err != nil {
		// try fallback
		id, err = readKenv()
	}
	if err != nil {
		return ""
	}
	return id
}

func readHostid() (string, error) {
	buf, err := readFile(hostidPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.Trim(string(buf), "\n")), nil
}

func readKenv() (string, error) {
	buf := &bytes.Buffer{}
	err := run(buf, os.Stderr, "kenv", "-q", "smbios.system.uuid")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.Trim(buf.String(), "\n")), nil
}
