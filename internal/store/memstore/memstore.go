package memstore

import (
	"github.com/Hasan-Al-Khazraji/CuckooCache/internal/store/lru"
)

// Thin wrapper around LRU to convert keys to strings and return eviction info

type Store struct {
	lru *lru.LRU
}

func New(capacity int) *Store {
	return &Store{
		lru: lru.New(capacity),
	}
}

func (s *Store) Get(key []byte) ([]byte, bool) {
	return s.lru.Get(string(key))
}

func (s *Store) Put(key, value []byte) (evictedKey, evictedVal []byte, evicted bool) {
	ek, ev, evicted := s.lru.Put(string(key), value)
	if !evicted {
		return nil, nil, false
	}
	return []byte(ek), ev, true
}

func (s *Store) Delete(key []byte) bool {
	return s.lru.Delete(string(key))
}

func (s *Store) Len() int {
	return s.lru.Len()
}
