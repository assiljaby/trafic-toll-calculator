package main

import "github.com/assiljaby/trafic-toll-calculator/types"

type Storer interface {
	insert(types.Distance) error
}

type MemoryStore struct {}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (ms *MemoryStore) insert(distance types.Distance) error {
	return nil
}