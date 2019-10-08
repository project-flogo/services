package cache

import (
	"sync"

	"github.com/project-flogo/services/flow-state/flow"
)

type memCache struct {
	steps     map[string][]*flow.StepInfo
	snapshots map[string]*flow.Snapshot
}

var lock = sync.RWMutex{}

func (c *memCache) AddSteps(key string, value *flow.StepInfo) {
	lock.Lock()
	if c.steps == nil {
		c.steps = make(map[string][]*flow.StepInfo)
	}
	c.steps[key] = append(c.steps[key], value)
	lock.Unlock()
}

func (c *memCache) AppendSteps(key string, value *flow.StepInfo) {
	lock.Lock()
	if c.steps == nil {
		c.steps = make(map[string][]*flow.StepInfo)
	}
	steps, ok := c.steps[key]
	if ok && len(steps) > 0 {
		steps = append(steps, value)
	} else {
		steps = []*flow.StepInfo{value}
	}

	c.steps[key] = steps
	lock.Unlock()
}

func (c *memCache) GetSteps(key string) []*flow.StepInfo {
	return c.steps[key]
}

func (c *memCache) DeleteSteps(key string) {
	lock.Lock()
	delete(c.steps, key)
	lock.Unlock()
}

func (c *memCache) AddSnapshots(key string, value *flow.Snapshot) {
	lock.Lock()
	if c.snapshots == nil {
		c.snapshots = make(map[string]*flow.Snapshot)
	}
	c.snapshots[key] = value
	lock.Unlock()
}

func (c *memCache) GetSnapshots(key string) *flow.Snapshot {
	return c.snapshots[key]
}

func (c *memCache) DeleteSnapshots(key string) {
	lock.Lock()
	delete(c.snapshots, key)
	lock.Unlock()
}

func (c *memCache) GetAllSnapshots() map[string]*flow.Snapshot {
	return c.snapshots
}

func NewCache() *memCache {
	return &memCache{
		steps:     make(map[string][]*flow.StepInfo),
		snapshots: make(map[string]*flow.Snapshot),
	}
}
