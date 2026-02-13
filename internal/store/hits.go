package store

import (
	"sync"
	"app/internal/model"
)

type HitsStore struct {
	mu   sync.RWMutex
	hits map[string]*model.Hits
}

func NewHitsStore() *HitsStore {
	return &HitsStore{
		hits: make(map[string]*model.Hits),
	}
}

func (s *HitsStore) GetHits(key string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	hits, exists := s.hits[key]
	if !exists {
		return 0, false
	}
	return int(hits.Total), exists
}

func (s *HitsStore) SetHits(key string, hits int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if existingHits, exists := s.hits[key]; exists {
		existingHits.Total = int64(hits)
	} else {
		s.hits[key] = &model.Hits{Total: int64(hits)}
	}
}

func (s *HitsStore) DeleteHits(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.hits, key)
}
