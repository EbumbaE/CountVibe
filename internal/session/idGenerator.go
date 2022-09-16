package session

import "sync"

type IDGenerator struct {
	sync.Mutex
	id int64
}

func (g *IDGenerator) newID() int64 {
	g.Lock()
	defer g.Unlock()
	g.id++
	return g.id
}
