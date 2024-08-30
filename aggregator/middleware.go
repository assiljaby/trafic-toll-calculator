package main

import (
	"time"

	"github.com/assiljaby/trafic-toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type MetricsMiddleware struct {
	errCounterAgg  prometheus.Counter
	errCounterCalc prometheus.Counter
	reqCounterAgg  prometheus.Counter
	reqCounterCalc prometheus.Counter
	reqLatencyCalc prometheus.Histogram
	reqLatencyAgg  prometheus.Histogram
	next           Aggregator
}

func NewMetricsMiddleware(next Aggregator) *MetricsMiddleware {
	errCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "aggregate",
	})
	errCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "calculate",
	})
	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregate",
	})
	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "calculate",
	})
	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "calculate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	return &MetricsMiddleware{
		next:           next,
		errCounterAgg:  errCounterAgg,
		errCounterCalc: errCounterCalc,
		reqCounterAgg:  reqCounterAgg,
		reqCounterCalc: reqCounterCalc,
		reqLatencyAgg:  reqLatencyAgg,
		reqLatencyCalc: reqLatencyCalc,
	}
}

func (m *MetricsMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		m.reqLatencyAgg.Observe(time.Since(start).Seconds())
		m.reqCounterAgg.Inc()
		if err != nil {
			m.errCounterAgg.Inc()
		}
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}

func (m *MetricsMiddleware) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		m.reqLatencyCalc.Observe(time.Since(start).Seconds())
		m.reqCounterCalc.Inc()
		if err != nil {
			m.errCounterCalc.Inc()
		}
	}(time.Now())
	inv, err = m.next.CalculateInvoice(obuID)
	return
}

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
