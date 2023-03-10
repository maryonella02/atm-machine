package cache

import (
	"container/list"
	"sync"
	"time"
)

type Cache struct {
	mu   sync.Mutex
	data map[string]*struct {
		key        string
		value      interface{}
		expiration time.Time
		element    *list.Element
	}
	list       *list.List
	expiration time.Duration
}

func NewLRUCache(expiration time.Duration) *Cache {
	return &Cache{
		data: make(map[string]*struct {
			key        string
			value      interface{}
			expiration time.Time
			element    *list.Element
		}),
		list:       list.New(),
		expiration: expiration,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.data[key]; ok {
		if time.Now().Before(e.expiration) {
			c.list.MoveToFront(e.element)
			return e.value, true
		}
		c.list.Remove(e.element)
		delete(c.data, key)
	}
	return nil, false
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.data[key]; ok {
		c.list.MoveToFront(e.element)
		e.value = value
		e.expiration = time.Now().Add(c.expiration)
	} else {
		e := &struct {
			key        string
			value      interface{}
			expiration time.Time
			element    *list.Element
		}{
			key:        key,
			value:      value,
			expiration: time.Now().Add(c.expiration),
			element:    c.list.PushFront(nil),
		}
		e.element.Value = e
		c.data[key] = e
		if c.list.Len() > 0 && c.expiration > 0 && c.list.Len() > len(c.data) {
			c.deleteOldest()
		}
	}
}

func (c *Cache) deleteOldest() {
	if ele := c.list.Back(); ele != nil {
		c.list.Remove(ele)
		delete(c.data, ele.Value.(*struct {
			key        string
			value      interface{}
			expiration time.Time
			element    *list.Element
		}).key)
	}
}
