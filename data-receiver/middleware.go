package main

import (
	"time"

	"github.com/assiljaby/trafic-toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type loggerMiddleware struct {
	next DataProducer
}

func newLoggerMiddleware(next DataProducer) *loggerMiddleware {
	return &loggerMiddleware{
		next: next,
	}
}

func (lm *loggerMiddleware) ProduceData(data types.OBUData) error {
	start := time.Now()
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuID": data.OBUID,
			"lat": data.Lat,
			"long": data.Long,
			"Took": time.Since(start),
		}).Info("Producing to Kafka")
	}(start)
	return lm.next.ProduceData(data)
}