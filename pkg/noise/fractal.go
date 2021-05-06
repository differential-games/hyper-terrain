package noise

import (
	"willbeason/hyper-terrain/pkg/vector"
)

const (
	scale    = 1.6180339887498948482045868343656381177203091798057628621354486227
	invScale = 1.0 / scale
	cosAngle = -0.362374890080480119958646637474986899360865544005598546450156788
	sinAngle = 0.9320324238132276215340316686911046913861543007830788182892528729

	offsetX = scale * float64(size)
	offsetY = (scale - 1) * (scale - 1) * float64(size)
)

type Transform struct {
	Offset vector.V2
}

type Fractal struct {
	Value
}

func (f *Fractal) Cubic(x, y float64) float64 {
	sum := 0.0
	c := 1.0
	for i := 0; i < 15; i++ {
		sum += f.Value.Cubic(x, y) * c
		x, y = x*cosAngle-y*sinAngle, x*sinAngle+y*cosAngle
		x, y = x+offsetX, y+offsetY
		x *= scale
		y *= scale
		c *= invScale
	}
	return sum
}

func (f *Fractal) Linear(x, y float64) float64 {
	sum := 0.0
	c := 1.0
	for i := 0; i < 15; i++ {
		sum += f.Value.Linear(x, y) * c
		x, y = x*cosAngle-y*sinAngle, x*sinAngle+y*cosAngle
		x, y = x+offsetX, y+offsetY
		x *= scale
		y *= scale
		c *= invScale
	}
	return sum
}
