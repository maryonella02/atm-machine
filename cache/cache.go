package cache

import (
	"container/list"
	"sync"
	"time"
)

type Cache struct {
	mu   sync.Mutex
	data map[string]*entry
	// list is a doubly-linked list of entries, used for eviction
	list       *list.List
	expiration time.Duration
}

type entry struct {
	key        string
	value      interface{}
	expiration time.Time
	element    *list.Element
}

// NewLRUCache creates a new LRU cache with the specified expiration time.
func NewLRUCache(expiration time.Duration) *Cache {
	return &Cache{
		data:       make(map[string]*entry),
		list:       list.New(),
		expiration: expiration,
	}
}

// Get returns the value for the given key, if it exists and has not expired.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.data[key]; ok {
		if time.Now().Before(e.expiration) {
			c.moveToFront(e)
			return e.value, true
		}
		c.delete(e)
	}
	return nil, false
}

// Set sets the value for the given key.
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.data[key]; ok {
		c.delete(e)
	}
	e := &entry{
		key:        key,
		value:      value,
		expiration: time.Now().Add(c.expiration),
		element:    c.list.PushFront(&entry{}),
	}
	*e.element.Value.(*entry) = *e
	c.data[key] = e
}

func (c *Cache) delete(e *entry) {
	delete(c.data, e.key)
	c.list.Remove(e.element)
}

func (c *Cache) moveToFront(e *entry) {
	c.list.MoveToFront(e.element)
}
