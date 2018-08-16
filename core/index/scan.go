package index

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/ossman11/sip/core/def"
)

type Scan struct {
	httpClient *http.Client
	Result     map[string]Connection
	Running    bool
}

func (i *Scan) endGo(c chan bool) {
	if c != nil {
		c <- true
	}
}

func (i *Scan) startChan(buf int) (int, chan bool) {
	if buf < 1 {
		buf = 1
	}
	return 0, make(chan bool, buf)
}

func (i *Scan) awaitChan(cc int, ch chan bool) {
	for cc > 0 {
		<-ch
		cc--
	}
}

func (i *Scan) getIP(ip net.IP, c chan bool) {
	// Prepare the target string
	str := ip.String() + ":1670" // 1670

	res, err := i.httpClient.Get("http://" + str + def.APIIndexJoin)
	if err == nil {
		var m Type
		dec := json.NewDecoder(res.Body)
		err = dec.Decode(&m)
		if err != nil {
			fmt.Println(err)
			return
		}

		id := ParseIP(ip)
		i.Result[id.String()] = Connection{
			IP:   ip,
			ID:   id,
			Type: m,
		}
	}

	i.endGo(c)
}

func (i *Scan) walkIP(ipnet *net.IPNet, c chan bool) {
	defer i.endGo(c)

	fmt.Println(ipnet)

	// Convert IP and mask into integers
	ip := binary.BigEndian.Uint32([]byte(ipnet.IP))
	srcIP := ip
	mask := binary.BigEndian.Uint32([]byte(ipnet.Mask))

	ip &= mask
	mask = mask | 0xffffff00
	// Skip networks that have to large masks
	if 0xffffffff-mask > 1<<8 {
		return
	}

	cc, ch := i.startChan(0)

	// Walk over all possible addresses that are on the IP range
	for mask < 0xffffffff {
		if ip != srcIP {
			// Construct new target IP
			var buf [4]byte
			binary.BigEndian.PutUint32(buf[:], ip)
			tar := net.IP(buf[:])

			// Process the target IP address
			go i.getIP(tar, ch)
			cc++
		}
		// Prepare the next ip address
		mask++
		ip++

		if (0xffffffff-mask)%100 == 0 {
			fmt.Println(0xffffffff - mask)
		}
	}

	i.awaitChan(cc, ch)
}

func (i *Scan) scanIP(ipnet *net.IPNet, c chan bool) bool {
	filter := func(ipnet *net.IPNet) bool {
		if ipnet.IP[0] == 127 {
			return false
		}
		return true
	}

	ip4 := ipnet.IP.To4()
	if ip4 != nil {
		addr := &net.IPNet{
			IP:   ip4,
			Mask: ipnet.Mask[len(ipnet.Mask)-4:],
		}
		if filter(addr) {
			go i.walkIP(addr, c)
		}
	} else {
		return false
	}
	return true
}

func (i *Scan) scanInteface(iface net.Interface, c chan bool) {

	as, err := iface.Addrs()
	if err != nil {
		log.Fatal(err)
	}

	cc, ch := i.startChan(len(as) + 1)

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if ok {
			if i.scanIP(ipnet, ch) {
				cc++
			}
		}
	}

	i.awaitChan(cc, ch)
	i.endGo(c)
}

func (i *Scan) Scan() {
	if i.Running {
		return
	}
	i.Running = true

	faces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	cc, ch := i.startChan(0)

	for _, v := range faces {
		if v.Flags&net.FlagLoopback != 0 || v.Flags&net.FlagUp == 0 {
			continue
		}
		go i.scanInteface(v, ch)
		cc++
	}

	i.awaitChan(cc, ch)
	i.Running = false
}

func NewScan() Scan {
	timeout := time.Duration(30 * time.Second)
	i := Scan{
		Result: map[string]Connection{},
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
	return i
}
