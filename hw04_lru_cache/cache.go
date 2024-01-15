package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type CacheItem struct {
	Key   Key
	Value interface{}
}

type lruCache struct {
	capacity int
	queue    List

	mx    *sync.Mutex
	items map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if c.capacity <= 0 {
		return false
	}

	c.mx.Lock()
	defer c.mx.Unlock()

	if i, ok := c.items[key]; ok {
		i.Value.(*CacheItem).Value = value
		c.queue.MoveToFront(i)

		return true
	}

	if c.capacity == c.queue.Len() {
		c.Clear()
	}

	item := NewCacheItem(key, value)

	c.items[key] = c.queue.PushFront(&item)

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	if i, ok := c.items[key]; ok {
		c.queue.MoveToFront(i)

		return i.Value.(*CacheItem).Value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	i := c.queue.Back()

	key := i.Value.(*CacheItem).Key

	delete(c.items, key)
	c.queue.Remove(i)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		mx:       &sync.Mutex{},
		items:    make(map[Key]*ListItem, capacity),
	}
}

func NewCacheItem(key Key, v interface{}) CacheItem {
	return CacheItem{
		Key:   key,
		Value: v,
	}
}
