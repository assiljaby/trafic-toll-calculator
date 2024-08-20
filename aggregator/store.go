package main

import "github.com/assiljaby/trafic-toll-calculator/types"

type Storer interface {
	insert(types.Distance) error
}

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (ms *MemoryStore) insert(distance types.Distance) error {
	ms.data[distance.OBUID] += distance.Value
	return nil
}
