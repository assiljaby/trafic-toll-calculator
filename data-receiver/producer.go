package main

import (
	"encoding/json"
	"fmt"

	"github.com/assiljaby/trafic-toll-calculator/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type DataProducer interface {
	ProduceData(types.OBUData) error
}

type kafkaProducer struct {
	producer *kafka.Producer
	topic string
}

func NewKafkaProducer(topic string) (DataProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}
	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return &kafkaProducer{
		producer: p,
		topic: topic,
	}, nil
}

func (kp *kafkaProducer) ProduceData(data types.OBUData) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.topic, Partition: kafka.PartitionAny},
		Value:          d,
	}, nil)
}