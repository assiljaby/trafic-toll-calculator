package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

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
	for {
		for _, id := range OBUIDs {
			lat, long := generateLatLong()
			data := OBUData{
				OBUID: id,
				Lat: lat,
				Long: long,
			}
			fmt.Printf("%v\n", data)
		}
		time.Sleep(sendInterval)
	}
}