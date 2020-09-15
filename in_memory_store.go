package main

// NewInMemoryPlayerStore will return a new store in memory
func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

// InMemoryPlayerStore is the db abstraction for in-memory
type InMemoryPlayerStore struct {
	store map[string]int
}

// RecordWin will take a name and append a +=1 to the user in the store dict
func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.store[name]++
}

// GetPlayerScore takes a name and returns that players score as an int
func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}
