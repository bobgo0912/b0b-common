package balancer

import (
	"github.com/bobgo0912/b0b-common/pkg/server/model"
	"hash/crc32"
	"math/rand"
	"sync"
	"time"
)

const Salt = "%#!"

//func init() {
//	factories[P2CBalancer] = NewP2C
//}

type host struct {
	*model.EtcdReg
	name string
	load uint64
}

// P2C refer to the power of 2 random choice
type P2C struct {
	sync.RWMutex
	hosts   []*host
	rnd     *rand.Rand
	loadMap map[string]*host
}

// NewP2C create new P2C balancer
func NewP2C() *P2C {
	p := &P2C{
		hosts:   []*host{},
		loadMap: make(map[string]*host),
		rnd:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return p
}

// Add new host to the balancer
func (p *P2C) Add(hostName *model.EtcdReg) {
	p.Lock()
	defer p.Unlock()
	if _, ok := p.loadMap[hostName.Key]; ok {
		return
	}

	h := &host{name: hostName.HostName, EtcdReg: hostName, load: 0}
	p.hosts = append(p.hosts, h)
	p.loadMap[hostName.Key] = h
}

// Remove new host from the balancer
func (p *P2C) Remove(key string) {
	p.Lock()
	defer p.Unlock()
	if _, ok := p.loadMap[key]; !ok {
		return
	}
	delete(p.loadMap, key)
	for i, h := range p.hosts {
		if h.name == key {
			p.hosts = append(p.hosts[:i], p.hosts[i+1:]...)
			return
		}
	}
}

// Balance selects a suitable host according to the key value
func (p *P2C) Balance(key string) (*model.EtcdReg, error) {
	p.RLock()
	defer p.RUnlock()

	if len(p.hosts) == 0 {
		return nil, NoHostError
	}

	n1, n2 := p.hash(key)
	host := n2
	if p.loadMap[n1.Key].load <= p.loadMap[n1.Key].load {
		host = n1
	}
	return host, nil
}

func (p *P2C) hash(key string) (*model.EtcdReg, *model.EtcdReg) {
	var n1, n2 *model.EtcdReg
	if len(key) > 0 {
		saltKey := key + Salt
		n1 = p.hosts[crc32.ChecksumIEEE([]byte(key))%uint32(len(p.hosts))].EtcdReg
		n2 = p.hosts[crc32.ChecksumIEEE([]byte(saltKey))%uint32(len(p.hosts))].EtcdReg
		return n1, n2
	}
	n1 = p.hosts[p.rnd.Intn(len(p.hosts))].EtcdReg
	n2 = p.hosts[p.rnd.Intn(len(p.hosts))].EtcdReg
	return n1, n2
}

// Inc refers to the number of connections to the server `+1`
func (p *P2C) Inc(host string) {
	p.Lock()
	defer p.Unlock()

	h, ok := p.loadMap[host]

	if !ok {
		return
	}
	h.load++
}

// Done refers to the number of connections to the server `-1`
func (p *P2C) Done(host string) {
	p.Lock()
	defer p.Unlock()

	h, ok := p.loadMap[host]

	if !ok {
		return
	}

	if h.load > 0 {
		h.load--
	}
}
