package vector

import (
	"math"
)

// V2 represents a point in 2D space.
type V2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Angle returns a unit vector making the passed angle with the x-axis.
func Angle(theta float64) V2 {
	return V2{
		X: math.Cos(theta),
		Y: math.Sin(theta),
	}
}
