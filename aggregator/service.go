package main

import (
	"fmt"

	"github.com/assiljaby/trafic-toll-calculator/types"
)

const basePrice = 3.15

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (ia *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	fmt.Println("Processing and storing distance:", distance)
	return ia.store.insert(distance)
}

func (ia *InvoiceAggregator) CalculateInvoice(obuID int) (*types.Invoice, error) {
	sumDistance, err := ia.store.get(obuID)
	if err != nil {
		return nil, err
	}
	return &types.Invoice{
		OBUID:         obuID,
		TotalDistance: sumDistance,
		TotalAmount:   sumDistance * basePrice,
	}, nil
}
