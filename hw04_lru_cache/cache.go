package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.items[key]; exists {
		c.items[key].Value = cacheItem{key: key, value: value}
		c.queue.MoveToFront(c.items[key])
		return true
	}

	if c.queue.Len() == c.capacity {
		i := c.queue.Back()
		iValue, ok := i.Value.(cacheItem)
		if !ok {
			return false
		}
		delete(c.items, iValue.key)
		c.queue.Remove(c.queue.Back())
	}

	c.items[key] = c.queue.PushFront(cacheItem{key: key, value: value})
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if i, exists := c.items[key]; exists {
		iValue, ok := i.Value.(cacheItem)
		if !ok {
			return nil, false
		}
		c.queue.MoveToFront(c.items[key])
		return iValue.value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.queue = NewList()
	for k := range c.items {
		delete(c.items, k)
	}
}
