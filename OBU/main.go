package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

const wsEndpoint = "ws://127.0.0.1:30000/ws"

var sendInterval = time.Second

type OBUData struct {
	OBUID int     `json:"obuID"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

func generateCoord() float64 {
	return float64(rand.Intn(100)) + 1 + (rand.Float64())
}

func generateLatLong() (float64, float64) {
	return generateCoord(), generateCoord()
}

func generateOBUIDs(n int) []int {
	ids := make([]int, n)

	for i := range ids {
		ids[i] = rand.Intn(math.MaxInt)
	}

	return ids
}

func main() {
	OBUIDs := generateOBUIDs(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for _, id := range OBUIDs {
			lat, long := generateLatLong()
			data := OBUData{
				OBUID: id,
				Lat: lat,
				Long: long,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(sendInterval)
	}
}