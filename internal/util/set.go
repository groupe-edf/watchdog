package util

import (
	"sync"

	"github.com/groupe-edf/watchdog/internal/core/models"
)

// Set struct
type Set struct {
	items []models.Issue
	sync.RWMutex
}

// ConcurrentSliceItem Concurrent slice item
type ConcurrentSliceItem struct {
	Index int
	Value interface{}
}

// NewSet create new thread safe set
func NewSet() *Set {
	return &Set{
		items: make([]models.Issue, 0),
	}
}

// Add add
func (s *Set) Add(items []models.Issue) {
	s.Lock()
	defer s.Unlock()
	s.items = append(s.items, items...)
}

// Clear removes all items from the set
func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.items = make([]models.Issue, 0)
}

// Len returns the number of items in a set.
func (s *Set) Len() int {
	return len(s.items)
}

// List returns a slice of all items
func (s *Set) List() []models.Issue {
	s.RLock()
	defer s.RUnlock()
	return s.items
}

// Iter Iterates over the items in the concurrent slice
func (s *Set) Iter() <-chan ConcurrentSliceItem {
	c := make(chan ConcurrentSliceItem)
	f := func() {
		s.Lock()
		defer s.Unlock()
		for index, value := range s.items {
			c <- ConcurrentSliceItem{index, value}
		}
		close(c)
	}
	go f()
	return c
}
