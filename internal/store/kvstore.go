package store

import (
	"sync"
	"time"
)

type KVStore struct {
	mu    sync.RWMutex
	items map[string]Item
}

func New(cleanupInternal time.Duration) *KVStore {
	s := &KVStore{
		items: make(map[string]Item),
	}

	if cleanupInternal > 0 {
		s.StartJanitor(cleanupInternal)
	}

	return s
}

func (s *KVStore) Set(key string, value any, ttl time.Duration) {
	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = Item{
		Value:      value,
		Expiration: expiration,
	}

	setTotalOps.Inc()
	activeKeys.Inc()
}

func (s *KVStore) Get(key string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, found := s.items[key]
	if !found || item.Expired() {
		return nil, false
	}

	getTotalOps.Inc()

	return item.Value, true
}

func (s *KVStore) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, found := s.items[key]
	if !found {
		return false
	}

	delete(s.items, key)

	deleteTotalOps.Inc()
	activeKeys.Dec()
	return true
}

func (s *KVStore) DeleteExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixNano()
	for k, v := range s.items {
		if v.Expiration > 0 && now > v.Expiration {
			delete(s.items, k)
			expiredKeysTotal.Inc()
			activeKeys.Dec()
		}
	}
}

func (s *KVStore) StartJanitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			s.DeleteExpired()
		}
	}()
}
