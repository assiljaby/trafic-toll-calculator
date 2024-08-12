package main

import (
	"math"

	"github.com/assiljaby/trafic-toll-calculator/types"
)

type CalculateServicer interface {
	CalculateDistance(types.OBUData) (float64, error)
}

type point struct {
	lat float64
	long float64
}

type CalculateService struct {
	points []point
}

func NewCalculateService() CalculateServicer {
	return &CalculateService{
		points: make([]point, 0),
	}
}

func (cs *CalculateService) CalculateDistance(data types.OBUData) (float64, error) {
	// fmt.Println("Calculating Distance...")
	distance := 0.0
	if nPoints := len(cs.points); nPoints > 1 {
		distance = calculateDistance(cs.points[nPoints-1].lat, cs.points[nPoints-1].long, data.Lat, data.Long)
	}
	// TODO: Implement a buffer ring, appending to a slice is not optimal
	cs.points = append(cs.points, point{data.Lat, data.Long})
	return distance, nil
}

func calculateDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow((x2 - x1), 2) + math.Pow((y2 - y1), 2))
}