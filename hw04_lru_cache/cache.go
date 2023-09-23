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
	mutex    sync.Mutex
	queue    List
	items    map[Key]*ListItem
}

type cacheValue struct {
	Key   Key
	Value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if item, ok := l.items[key]; ok {
		item.Value = cacheValue{Key: key, Value: value}
		l.queue.MoveToFront(item)
		l.items[key] = item
		return true
	}

	item := l.queue.PushFront(cacheValue{
		Key:   key,
		Value: value,
	})
	if l.queue.Len() > l.capacity {
		oldest := l.queue.Back()
		l.queue.Remove(oldest)

		if cacheVal, ok := oldest.Value.(cacheValue); ok {
			delete(l.items, cacheVal.Key)
		}
	}

	l.items[key] = item
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	item, ex := l.items[key]
	if !ex {
		return nil, false
	}

	if cacheVal, ok := item.Value.(cacheValue); ok {
		l.queue.MoveToFront(item)
		return cacheVal.Value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
