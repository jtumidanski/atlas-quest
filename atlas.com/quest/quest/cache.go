package quest

import (
	"sync"
)

type cache struct {
	quests map[uint16]Model
	lock   sync.RWMutex
}

var c *cache
var once sync.Once

func GetCache() *cache {
	once.Do(func() {
		c = &cache{
			quests: make(map[uint16]Model, 0),
			lock:   sync.RWMutex{},
		}
	})
	return c
}

func (c *cache) Init() error {
	quests, err := readQuests()
	if err != nil {
		return err
	}

	c.lock.Lock()
	for _, q := range quests {
		c.quests[q.Id()] = q
	}
	c.lock.Unlock()
	return nil
}

//func (c *cache) GetFile(id uint32) (*Model, error) {
//	c.lock.RLock()
//	if val, ok := c.quests[id]; ok {
//		c.lock.RUnlock()
//		return &val, nil
//	} else {
//		c.lock.RUnlock()
//		c.lock.Lock()
//		s, err := readStatistics(id)
//		if err != nil {
//			c.lock.Unlock()
//			return nil, err
//		}
//		c.quests[id] = *s
//		c.lock.Unlock()
//		return s, nil
//	}
//}
