package main

import (
	"log"

	"github.com/assiljaby/trafic-toll-calculator/aggregator/client"
)

const topic = "obudata"
const aggEndpoint = "http://127.0.0.1:3000"

func main() {
	var (
		err  error
		csvc CalculateServicer
	)

	csvc = NewCalculateService()
	csvc = NewLoggerMiddleware(csvc)

	httpClient := client.NewHTTPClient(aggEndpoint)
	// grpcClient, err := client.NewGRPCClient(aggEndpoint)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	kConsumer, err := NewKafkaConsumer(topic, csvc, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	kConsumer.Start()
}
