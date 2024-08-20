package main

import (
	"log"

	"github.com/assiljaby/trafic-toll-calculator/aggregator/client"
)

const topic = "obudata"
const aggEndpoint = "http://127.0.0.1:3000/aggregate"

func main() {
	var (
		err  error
		csvc CalculateServicer
	)

	csvc = NewCalculateService()
	csvc = NewLoggerMiddleware(csvc)

	kConsumer, err := NewKafkaConsumer(topic, csvc, client.NewClient(aggEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	kConsumer.Start()
}
