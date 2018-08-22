package index

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
)

var HWIDCache string = ""

func HWID() string {
	if HWIDCache != "" {
		return HWIDCache
	}

	ret := fetchHWID()
	if ret == "" {
		fmt.Println("Failed to fetch HWID of the system!")
	}

	faces, err := net.Interfaces()
	if err == nil {
		for _, v := range faces {
			as, err := v.Addrs()
			if err != nil {
				continue
			}
			for _, a := range as {
				ipnet, ok := a.(*net.IPNet)
				if !ok {
					continue
				}
				ip4 := ipnet.IP.To4()
				if ip4 != nil {
					ret += "\nip:" + ip4.String()
				}
			}
		}
	}

	HWIDCache = ret
	return HWIDCache
}

func run(stdout, stderr io.Writer, cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdin = os.Stdin
	c.Stdout = stdout
	c.Stderr = stderr
	return c.Run()
}
