package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/assiljaby/trafic-toll-calculator/types"
	"github.com/gorilla/websocket"
)

func main() {
	rcvr, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", rcvr.handleWS)
	http.ListenAndServe(":30000", nil)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
type DataReceiver struct {
	msgch chan types.OBUData
	conn *websocket.Conn
	producer DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p DataProducer
		err error
		kafkaTopic = "obudata"
	)
	p, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}

	p = newLoggerMiddleware(p)

	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
		producer: p,
	}, nil
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReceiverLoop()
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.producer.ProduceData(data)
}

func (dr *DataReceiver) wsReceiverLoop() {
	fmt.Println("New OBU Connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("Read Error: ", err)
			continue
		}

		if err := dr.produceData(data); err != nil {
			fmt.Println("Kafka Produce Error: ", err)
		}
	}
}