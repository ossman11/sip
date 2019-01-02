package index

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

type ID string

func ParseHWID() ID {
	b := sha256.Sum256([]byte(HWID()))
	return ID(hex.EncodeToString(b[:]))
}

/* Unused
func ParseIP(ip net.IP) ID {
	b := sha256.Sum256([]byte(ip.String()))
	return ID(hex.EncodeToString(b[:]))
}

func ParseStr(str string) ID {
	return ID(str)
}
*/

func ParseByte(b []byte) ID {
	return ID(hex.EncodeToString(b[:]))
}

func NewIDConnection(n []*Node) ID {
	sort.Slice(n, func(i, j int) bool {
		return n[i].ID < n[j].ID
	})

	newID := ""

	for _, v := range n {
		newID += string(v.ID) + ":" + v.IP.String() + "\n"
	}

	b := sha256.Sum256([]byte(newID))
	return ParseByte(b[:])
}

func NewIDRoute(id []ID) ID {
	// Ensure that at least a single ID is provided
	if len(id) < 1 {
		return NewIDHash("")
	}

	newID := string(id[0]) + "\n" + string(id[len(id)-1])

	return NewIDHash(newID)
}

func NewIDHash(newID string) ID {
	b := sha256.Sum256([]byte(newID))
	return ParseByte(b[:])
}

func NewID(newID string) ID {
	b := []byte(newID)
	return ParseByte(b[:])
}

/* Unused
func (i *ID) Hash() []byte {
	ret, _ := hex.DecodeString(i.String())
	return ret
}

func (i *ID) String() string {
	return string(*i)
}

func (i *ID) Out(ip net.IP) ID {
	h := i.Hash()
	ipID := ParseIP(ip)
	d := ipID.Hash()
	r := i.Hash()

	for y := range h {
		r[y] = h[y] + d[y]
	}

	return ParseByte(r)
}

func (i *ID) In(ip net.IP) ID {
	h := i.Hash()
	ipID := ParseIP(ip)
	d := ipID.Hash()
	r := i.Hash()

	for y := range h {
		r[y] = byte(h[y] - d[y])
	}

	return ParseByte(r)
}
*/
