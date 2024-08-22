package main

import (
	"fmt"

	"github.com/assiljaby/trafic-toll-calculator/types"
)

type Storer interface {
	insert(types.Distance) error
	get(int) (float64, error)
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

func (ms *MemoryStore) get(obuID int) (float64, error) {
	distance, ok := ms.data[obuID]
	if !ok {
		return 0.0, fmt.Errorf("could not get the total didtance for this id: %d", obuID)
	}

	return distance, nil
}
