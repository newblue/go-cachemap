package cachemap

import "sync"

type Retriever func(string) (interface{}, bool)

type Cache struct {
	lock *sync.RWMutex
	data map[string]interface{}
}

func New() (c *Cache) {
	c = new(Cache)
	c.lock = new(sync.RWMutex)
	c.data = make(map[string]interface{}, 10)
	return
}

func (c *Cache) Get(key string, f Retriever) (v interface{}, ok bool) {
	if v, ok = c.get(key); !ok {
		if v, ok = f(key); ok {
			go c.set(key, v)
		}
	}
	return
}

func (c *Cache) Stale(key string) {
	go c.unset(key)
}

func (c *Cache) get(key string) (v interface{}, ok bool) {
	c.lock.RLock()
	v, ok = c.data[key]
	c.lock.RUnlock()
	return
}

func (c *Cache) set(key string, v interface{}) {
	c.lock.Lock()
	c.data[key] = v
	c.lock.Unlock()
}

func (c *Cache) unset(key string) {
	c.lock.Lock()
	c.data[key] = nil,false
	c.lock.Unlock()
}
