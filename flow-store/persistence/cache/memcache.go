package cache

import (
	"github.com/project-flogo/services/flow-store/flow"
	"sync"
)

type memCache struct {
	flows map[string]*flow.Flow
}

var lock = sync.RWMutex{}

func (c *memCache) AddFlow(key string, value *flow.Flow) {
	lock.Lock()
	if c.flows == nil {
		c.flows = make(map[string]*flow.Flow)
	}
	c.flows[key] = value
	lock.Unlock()
}

func (c *memCache) AllFlows() map[string]*flow.Flow {
	return c.flows
}

func (c *memCache) GetFlow(key string) *flow.Flow {
	return c.flows[key]
}

func (c *memCache) DeleteFlow(key string) {
	lock.Lock()
	delete(c.flows, key)
	lock.Unlock()
}

func NewCache() *memCache {
	return &memCache{
		flows: make(map[string]*flow.Flow),
	}
}
