package noise

import (
	"math"
	"math/rand"

	"willbeason/hyper-terrain/pkg/fixed"
)

// Value implements linearly interpolated value noise.
type Value struct {
	// noise is an array of the underlying noise.
	// Row (x) is the [0,size) bits of the index.
	// Column (y) is the [size,2*size) bits of the index.
	noise [size2]float64

	// Surprisingly, it is slower to precompute partial derivatives to reference
	// later than to just compute them on-the-fly.
}

// Fill generates the underlying noise which Value will interpolate.
//
// src is the source of randomness to use to generate noise.
func (v *Value) Fill(rnd *rand.Rand) {
	for i := 0; i < size2; i++ {
		v.noise[i] = rnd.Float64()
	}
}

// Linear returns the bilinear-interpolated value noise at the specified coordinates.
//
// Guarantees monotonic behavior between integral values.
// Guarantees behavior at (x, y) is equivalent to (x mod size, y mod size)
func (v *Value) Linear(x, y float64) float64 {
	// (x1, y1) is the bottom left of the cell containing (x, y).
	xi, xr := math.Modf(x)
	yi, yr := math.Modf(y)

	x1 := int(xi)
	if xi < 0 || xr < 0 {
		x1 = x1 - 1
		xr = 1 + xr
	}
	x1 = x1 & intMask

	y1 := int(yi)
	if yi < 0 || yr < 0 {
		y1 = y1 - 1
		yr = 1 + yr
	}
	y1 = y1 & intMask
	y1 = y1 << shift

	// Get the value at each corner surrounding the position.
	// The compiler optimizes away these assignments; this is for readability.
	//
	// The additions and bitwise-anding are offsets and moduli respectively.
	// Measured faster inlined rather than as a function or stored.
	vBottomLeft := v.noise[y1+x1]
	vBottomRight := v.noise[y1+((x1+1)&intMask)]
	vUpperLeft := v.noise[((y1+size)&int2Mask)+x1]
	vUpperRight := v.noise[((y1+size)&int2Mask)+((x1+1)&intMask)]

	// Linearly interpolate based on the four corners of the enclosing square.
	// Measured faster to store xryr as to
	// 1) try to eliminate the second use, or
	// 2) not store the value.
	xryr := xr * yr
	return xryr * vUpperRight +
		(yr - xryr) * vUpperLeft +
		(xr - xryr) * vBottomRight +
		(1.0 + xryr - xr - yr) * vBottomLeft
}

// LinearFloat is a convenience method for getting the bilinear-interpolated value noise at a pair
// of float64 coordinates.
func (v *Value) LinearFloat(x, y float64) float64 {
	return v.Linear(x, y)
}

// Cubic returns the bicubic-interpolated value noise at the specified coordinates.
func (v *Value) Cubic(x, y float64) float64 {
	// This consistently beats a Value noise generator which only uses floats as getting the
	// corresponding indices ends up costing several extra nanoseconds per call.

	// (x1, y1) is the bottom left of the cell containing (x, y).
	xi, xr := math.Modf(x)
	yi, yr := math.Modf(y)

	x1 := int(xi)
	if xi < 0 || xr < 0 {
		x1 = x1 - 1
		xr = 1 + xr
	}
	x1 = x1 & intMask

	y1 := int(yi)
	if yi < 0 || yr < 0 {
		y1 = y1 - 1
		yr = 1 + yr
	}
	y1 = y1 & intMask
	y1 = y1 << shift

	// (x0, y0) is the bottom left of the cell south-west of the cell containing (x, y).
	x0 := (x1 - 1) & intMask
	y0 := (y1 - size) & int2Mask

	// (x2, y2) is the top right of the cell containing (x, y).
	x2 := (x1 + 1) & intMask
	y2 := (y1 + size) & int2Mask

	// (x3, y3) is the top right of the cell north-east of the cell containing (x, y).
	x3 := (x2 + 1) & intMask
	y3 := (y2 + size) & int2Mask

	// Get the random noise in a 4x4 grid centered on (x, y).
	f00 := v.noise[x0+y0]
	f01 := v.noise[x0+y1]
	f02 := v.noise[x0+y2]
	f03 := v.noise[x0+y3]
	f10 := v.noise[x1+y0]
	f11 := v.noise[x1+y1]
	f12 := v.noise[x1+y2]
	f13 := v.noise[x1+y3]
	f20 := v.noise[x2+y0]
	f21 := v.noise[x2+y1]
	f22 := v.noise[x2+y2]
	f23 := v.noise[x2+y3]
	f30 := v.noise[x3+y0]
	f31 := v.noise[x3+y1]
	f32 := v.noise[x3+y2]
	f33 := v.noise[x3+y3]

	// Calculate the y partial derivatives over the grid.
	fy01 := (f02 - f00) / 2
	fy02 := (f03 - f01) / 2
	fy11 := (f12 - f10) / 2
	fy12 := (f13 - f11) / 2
	fy21 := (f22 - f20) / 2
	fy22 := (f23 - f21) / 2
	fy31 := (f32 - f30) / 2
	fy32 := (f33 - f31) / 2

	// Even though these are single-use, the compiler automatically optimizes away these assignments.
	// They're split out for readability.
	//
	// Calculate the x partial derivatives over the grid.
	fx11 := (f21 - f01) / 2
	fx12 := (f22 - f02) / 2
	fx21 := (f31 - f11) / 2
	fx22 := (f32 - f12) / 2
	// Calculate the mixed xy partial derivatives over the grid.
	fxy11 := (fy21 - fy01) / 2
	fxy12 := (fy22 - fy02) / 2
	fxy21 := (fy31 - fy11) / 2
	fxy22 := (fy32 - fy12) / 2

	// It is slower to precompute squares.
	return fixed.V4{
		1, xr, xr * xr, xr * xr * xr,
	}.TimesM1().TimesMatrix(fixed.Matrix{
		f11, f12, fy11, fy12,
		f21, f22, fy21, fy22,
		fx11, fx12, fxy11, fxy12,
		fx21, fx22, fxy21, fxy22,
	}).TimesM2().Dot(fixed.V4{
		1, yr, yr * yr, yr * yr * yr,
	})
}

// CubicFloat is a convenience method for getting the bicubic-interpolated value noise at a pair of
// float64 coordinates.
func (v *Value) CubicFloat(x, y float64) float64 {
	return v.Cubic(x, y)
}
