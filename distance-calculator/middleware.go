package main

import (
	"time"

	"github.com/assiljaby/trafic-toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type LoggerMiddleware struct {
	next CalculateServicer
}

func NewLoggerMiddleware(next CalculateServicer) CalculateServicer {
	return LoggerMiddleware{
		next: next,
	}
} 

func (lm LoggerMiddleware) CalculateDistance(data types.OBUData) (dist float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err": err,
			"distance": dist,
		}).Info("Calculating Distance:")
	}(time.Now())
	dist, err = lm.next.CalculateDistance(data)
	return
}