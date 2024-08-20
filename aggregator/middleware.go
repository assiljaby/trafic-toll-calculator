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
			"OBUID":    distance.OBUID,
			"unix":     distance.Unix,
		}).Info("AggregateDistance")
	}(time.Now())

	err = l.next.AggregateDistance(distance)
	return
}
