package main

import "sync"

// NewInMemoryPlayerStore will return a new store in memory
func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		map[string]int{},
		sync.RWMutex{},
	}
}

// InMemoryPlayerStore is the db abstraction for in-memory
type InMemoryPlayerStore struct {
	store map[string]int
	lock  sync.RWMutex
}

// RecordWin will take a name and append a +=1 to the user in the store dict
func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.store[name]++
}

// GetPlayerScore takes a name and returns that players score as an int
func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.store[name]
}
