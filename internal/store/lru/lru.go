package lru

import (
	"container/list"
	"sync"
)

type entry struct {
	key   string
	value []byte
}

type LRU struct {
	mu  sync.Mutex // TODO: Consider using RWMutex for read-heavy workloads
	cap int
	ll  *list.List
	idx map[string]*list.Element
}

func New(capacity int) *LRU {
	if capacity <= 0 {
		panic("lru: capacity must be > 0")
	}
	return &LRU{
		cap: capacity,
		ll:  list.New(),
		idx: make(map[string]*list.Element),
	}
}

func (l *LRU) Get(k string) ([]byte, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if el, ok := l.idx[k]; ok {
		l.ll.MoveToFront(el)
		return el.Value.(*entry).value, true
	}
	return nil, false
}

func (l *LRU) Put(k string, v []byte) (evictedKey string, evictedVal []byte, evicted bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if el, ok := l.idx[k]; ok {
		en := el.Value.(*entry)
		en.value = v
		l.ll.MoveToFront(el)
		return "", nil, false
	}

	en := &entry{key: k, value: v}
	el := l.ll.PushFront(en)
	l.idx[k] = el

	if l.ll.Len() > l.cap {
		tail := l.ll.Back()
		l.ll.Remove(tail)
		ten := tail.Value.(*entry)
		delete(l.idx, ten.key)
		return ten.key, ten.value, true
	}

	return "", nil, false
}

func (l *LRU) Delete(k string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if el, ok := l.idx[k]; ok {
		l.ll.Remove(el)
		delete(l.idx, k)
		return true
	}
	return false
}

func (l *LRU) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.ll.Len()
}
