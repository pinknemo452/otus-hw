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
	m        sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.m.Lock()
	defer c.m.Unlock()

	if i, ok := c.items[key]; ok {
		i.Value.(*cacheItem).value = value
		c.queue.MoveToFront(i)
		return true
	}

	if c.queue.Len() == c.capacity {
		lastElement := c.queue.Back()
		item := lastElement.Value.(*cacheItem)
		delete(c.items, item.key)
		c.queue.Remove(c.queue.Back())
	}

	newItem := &cacheItem{key: key, value: value}
	c.items[key] = c.queue.PushFront(newItem)
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	if i, ok := c.items[key]; ok {
		c.queue.MoveToFront(i)
		return i.Value.(*cacheItem).value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.m.Lock()
	defer c.m.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
