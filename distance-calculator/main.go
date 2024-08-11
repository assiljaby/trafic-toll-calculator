package main

import "log"

const topic = "obudata"

func main() {
	var (
		err error
		csvc CalculateServicer
	)

	csvc = NewCalculateService()

	kConsumer, err := NewKafkaConsumer(topic, csvc)
	if err != nil {
		log.Fatal(err)
	}

	kConsumer.Start()
}