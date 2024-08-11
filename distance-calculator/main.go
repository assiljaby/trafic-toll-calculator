package main

import "log"

const topic = "obudata"

func main() {
	kConsumer, err := NewKafkaConsumer(topic)
	if err != nil {
		log.Fatal(err)
	}

	kConsumer.Start()
}