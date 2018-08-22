// +build linux

package index

import (
	"io/ioutil"
	"strings"
)

// Credit: https://github.com/denisbrodbeck/machineid/blob/master/id_linux.go

const (
	// dbusPath is the default path for dbus machine id.
	dbusPath = "/var/lib/dbus/machine-id"
	// dbusPathEtc is the default path for dbus machine id located in /etc.
	// Some systems (like Fedora 20) only know this path.
	// Sometimes it's the other way round.
	dbusPathEtc = "/etc/machine-id"
)

// machineID returns the uuid specified at `/var/lib/dbus/machine-id` or `/etc/machine-id`.
// If there is an error reading the files an empty string is returned.
// See https://unix.stackexchange.com/questions/144812/generate-consistent-machine-unique-id
func fetchHWID() string {
	id, err := ioutil.ReadFile(dbusPath)
	if err != nil {
		// try fallback path
		id, err = ioutil.ReadFile(dbusPathEtc)
	}
	if err != nil {
		return ""
	}
	return strings.TrimSpace(strings.Trim(string(id), "\n"))
}
