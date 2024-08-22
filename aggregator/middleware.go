package main

import (
	"time"

	"github.com/assiljaby/trafic-toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type LoggerMiddleware struct {
	next Aggregator
}

func NewLoggerMiddleware(next Aggregator) Aggregator {
	return &LoggerMiddleware{
		next: next,
	}
}

func (l *LoggerMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"Took":     time.Since(start),
			"err":      err,
			"distance": distance.Value,
			"obuid":    distance.OBUID,
			"unix":     distance.Unix,
		}).Info("AggregateDistance")
	}(time.Now())

	err = l.next.AggregateDistance(distance)
	return
}

func (l *LoggerMiddleware) CalculateInvoice(obuID int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
		)

		if invoice != nil {
			distance = invoice.TotalDistance
			amount = invoice.TotalAmount
		}

		logrus.WithFields(logrus.Fields{
			"took":          time.Since(start),
			"err":           err,
			"obuid":         obuID,
			"totalDistance": distance,
			"totalAmount":   amount,
		}).Info("Calculate Invoice")
	}(time.Now())

	invoice, err = l.next.CalculateInvoice(obuID)
	return
}
