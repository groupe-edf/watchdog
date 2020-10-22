package util

import (
	"sync"

	"github.com/groupe-edf/watchdog/internal/issue"
)

// Set struct
type Set struct {
	items map[issue.Issue]struct{}
	sync.RWMutex
}

// NewSet create new thread safe set
func NewSet() *Set {
	return &Set{
		items: make(map[issue.Issue]struct{}),
	}
}

// Add add
func (s *Set) Add(items []issue.Issue) {
	s.Lock()
	defer s.Unlock()
	for _, item := range items {
		s.items[item] = struct{}{}
	}
}

// Clear removes all items from the set
func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.items = make(map[issue.Issue]struct{})
}

// Has looks for the existence of an item
func (s *Set) Has(item issue.Issue) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.items[item]
	return ok
}

// IsEmpty checks for emptiness
func (s *Set) IsEmpty() bool {
	return s.Len() == 0
}

// Len returns the number of items in a set.
func (s *Set) Len() int {
	return len(s.List())
}

// List returns a slice of all items
func (s *Set) List() []issue.Issue {
	s.RLock()
	defer s.RUnlock()
	list := make([]issue.Issue, 0)
	for item := range s.items {
		list = append(list, item)
	}
	return list
}

// Remove deletes the specified item from the map
func (s *Set) Remove(item issue.Issue) {
	s.Lock()
	defer s.Unlock()
	delete(s.items, item)
}
