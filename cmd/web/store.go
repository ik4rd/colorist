package main

import (
	"sync"

	"github.com/ik4rd/colorist/internal/colormap"
)

type store struct {
	mu    sync.Mutex
	max   int
	items map[string]*colormap.Pixels
	order []string
}

func newStore(max int) *store {
	if max < 1 {
		max = 1
	}
	return &store{max: max, items: make(map[string]*colormap.Pixels)}
}

func (s *store) put(px *colormap.Pixels) string {
	id := newID()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[id] = px
	s.order = append(s.order, id)

	for len(s.order) > s.max {
		oldest := s.order[0]
		s.order = s.order[1:]
		delete(s.items, oldest)
	}

	return id
}

func (s *store) get(id string) (*colormap.Pixels, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	px, ok := s.items[id]

	return px, ok
}
