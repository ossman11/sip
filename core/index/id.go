package index

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
)

type ID [sha256.Size]byte

func ParseIP(ip net.IP) ID {
	return sha256.Sum256([]byte(ip.String()))
}

func ParseStr(str string) ID {
	var ret ID
	res, _ := hex.DecodeString(str)
	for i := range res {
		ret[i] = res[i]
	}
	return ret
}

func (i *ID) String() string {
	return hex.EncodeToString(i[:])
}

func (i *ID) Out(ip net.IP) ID {
	d := ParseIP(ip)
	r := ID{}

	for y := range i {
		r[y] = i[y] + d[y]
	}

	return r
}

func (i *ID) In(ip net.IP) ID {
	d := ParseIP(ip)
	r := ID{}

	for y := range i {
		r[y] = byte(i[y] - d[y])
	}

	return r
}
