package main

import (
	"encoding/json"

	"github.com/assiljaby/trafic-toll-calculator/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	isRunning bool
	calculatorSvc CalculateServicer
}

func NewKafkaConsumer(topic string, csvc CalculateServicer) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	if err = c.SubscribeTopics([]string{topic}, nil) ;err != nil {
		return nil, err
	}	 

	return &KafkaConsumer{
		consumer: c,
		isRunning: false,
		calculatorSvc: csvc,
	}, nil
}

func (kc *KafkaConsumer) Start() {
	logrus.Info("kafka started consuming")
	kc.isRunning = true
	kc.readMessageloop()
}

func (kc *KafkaConsumer) readMessageloop() {
	for kc.isRunning {
		msg, err := kc.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error: %s\n", err)
			continue
		}

		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("Json parsing error: %s\n", err)
			continue
		}

		distance, err := kc.calculatorSvc.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("Failed to calculate distance: %s\n", err)
			continue
		}
		// fmt.Printf("Distance: %.2f\n", distance)
		_ = distance
	}
}