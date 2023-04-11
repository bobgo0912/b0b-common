package balancer

import (
	"github.com/bobgo0912/b0b-common/pkg/server/model"
	"sync"
)

type P2cBalancer struct {
	L sync.Mutex
	M map[string]*P2C
}

func NewP2cBalancer() *P2cBalancer {
	return &P2cBalancer{L: sync.Mutex{}, M: map[string]*P2C{}}
}
func (p *P2cBalancer) Add(key string, reg *model.EtcdReg) {
	p.L.Lock()
	defer p.L.Unlock()
	c, ok := p.M[key]
	if !ok {
		p2C := NewP2C()
		p2C.Add(reg)
		p.M[key] = p2C
	} else {
		c.Add(reg)
	}
}

func (p *P2cBalancer) Remove(key string) {
	p.L.Lock()
	defer p.L.Unlock()
	c, ok := p.M[key]
	if !ok {
		return
	} else {
		c.Remove(key)
	}
}

func (p *P2cBalancer) Balance(key string) (*model.EtcdReg, error) {
	p.L.Lock()
	defer p.L.Unlock()
	c, ok := p.M[key]
	if !ok {
		return nil, NoHostError
	} else {
		return c.Balance(key)
	}
}
