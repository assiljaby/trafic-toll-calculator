package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/assiljaby/trafic-toll-calculator/types"
	"github.com/gorilla/websocket"
)

func main() {
	rcvr := NewDataReceiver()
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
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
	}
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReceiverLoop()
}

func (dr *DataReceiver) wsReceiverLoop() {
	fmt.Println("New OBU Connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("Read Error: ", err)
			continue
		}
		
		fmt.Printf("Received OBU data from <%d> :: <lat: %.2f, long: %.2f>\n", data.OBUID, data.Lat, data.Long)
		dr.msgch <- data
	}
}