package main

import (
	"encoding/json"
	"time"

	"github.com/assiljaby/trafic-toll-calculator/aggregator/client"
	"github.com/assiljaby/trafic-toll-calculator/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer      *kafka.Consumer
	isRunning     bool
	calculatorSvc CalculateServicer
	aggClient     *client.HTTPClient
}

func NewKafkaConsumer(topic string, csvc CalculateServicer, aggClient *client.HTTPClient) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	if err = c.SubscribeTopics([]string{topic}, nil); err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer:      c,
		isRunning:     false,
		calculatorSvc: csvc,
		aggClient:     aggClient,
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

		req := types.Distance{
			Value: distance,
			Unix:  time.Now().UnixNano(),
			OBUID: data.OBUID,
		}

		if err := kc.aggClient.AggregateInvoice(req); err != nil {
			logrus.Errorf("failed to send data to aggregate client: %v", err)
			continue
		}
	}
}
