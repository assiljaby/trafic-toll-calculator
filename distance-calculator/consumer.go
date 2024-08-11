package main

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	isRunning bool
}

func NewKafkaConsumer(topic string) (*KafkaConsumer, error) {
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
		fmt.Println(msg)
	}
}