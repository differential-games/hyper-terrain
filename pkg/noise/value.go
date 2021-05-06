package noise

import (
	"math/rand"

	"willbeason/hyper-terrain/pkg/fixed"
)

// Value implements linearly interpolated value noise.
type Value struct {
	// noise is an array of the underlying noise.
	// Row (x) is the [0,size) bits of the index.
	// Column (y) is the [size,2*size) bits of the index.
	noise [size2]fixed.F16
}

// Fill generates the underlying noise which Value will interpolate.
//
// src is the source of randomness to use to generate noise.
func (v *Value) Fill(src rand.Source) {
	for i := 0; i < size2; i++ {
		v.noise[i] = fixed.F16(src.Int63()).Remainder() // 0.0 to 1.0 - 2^-16
	}
}

// Nearest returns the noise nearest to (x, y).
func (v *Value) Nearest(x, y fixed.F16) fixed.F16 {
	// Take the modulus of the integral parts of each coordinate.
	// Each measured faster stored rather than recomputed 4 times.
	xi := x.Int() & intMask
	if x.Remainder() > fixed.Half16 {
		// Increment if closer to next.
		xi = (xi + 1) & intMask
	}
	yi := int(y>>revShift) & int2Mask
	if y.Remainder() > fixed.Half16 {
		// Increment if closer to next.
		yi = (yi + size) & int2Mask
	}

	return v.noise[yi+xi]
}

// Linear linearly interpolates noise.
//
// Guarantees monotonic behavior between integral values.
// Guarantees behavior at (x, y) is equivalent to (x mod size, y mod size)
func (v *Value) Linear(x, y fixed.F16) fixed.F32 {
	// Take the modulus of the integral parts of each coordinate.
	// Measured faster to store rather than recompute 4 times.
	xi := x.Int() & intMask
	yi := int(y>>revShift) & int2Mask

	// Get the value at each corner surrounding the position.
	// The compiler optimizes away these assignments; this is for readability.
	//
	// The additions and bitwise-anding are offsets and moduli respectively.
	// Measured faster inlined rather than as a function or stored.
	vBottomLeft := v.noise[yi+xi]
	vBottomRight := v.noise[yi+((xi+1)&intMask)]
	vUpperLeft := v.noise[((yi+size)&int2Mask)+xi]
	vUpperRight := v.noise[((yi+size)&int2Mask)+((xi+1)&intMask)]

	// Linearly interpolate based on the four corners of the enclosing square.

	// Measured faster to store these rather than recompute.
	xr := x.Remainder()
	yr := y.Remainder()
	// Measured faster to store xryr as to
	// 1) try to eliminate the second use, or
	// 2) not store the value.
	xryr := xr.Times(yr).F16()
	return xryr.Times(vUpperRight) +
		(yr - xryr).Times(vUpperLeft) +
		(xr - xryr).Times(vBottomRight) +
		(fixed.One16 + xryr - xr - yr).Times(vBottomLeft)
}

func (v *Value) LinearFloat(x, y float64) float64 {
	xf, yf := fixed.Float(x), fixed.Float(y)
	return v.Linear(xf, yf).Float64()
}

func (v *Value) Cubic(x, y fixed.F16) float64 {
	// This consistently beats a Value noise generator which only uses floats as getting the
	// corresponding indices ends up costing several extra nanoseconds per call.

	// (x1, y1) is the bottom left of the cell containing (x, y).
	x1 := x.Int() & intMask
	y1 := int(y>>revShift) & int2Mask

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
	f00 := v.noise[x0+y0].Float64()
	f01 := v.noise[x0+y1].Float64()
	f02 := v.noise[x0+y2].Float64()
	f03 := v.noise[x0+y3].Float64()
	f10 := v.noise[x1+y0].Float64()
	f11 := v.noise[x1+y1].Float64()
	f12 := v.noise[x1+y2].Float64()
	f13 := v.noise[x1+y3].Float64()
	f20 := v.noise[x2+y0].Float64()
	f21 := v.noise[x2+y1].Float64()
	f22 := v.noise[x2+y2].Float64()
	f23 := v.noise[x2+y3].Float64()
	f30 := v.noise[x3+y0].Float64()
	f31 := v.noise[x3+y1].Float64()
	f32 := v.noise[x3+y2].Float64()
	f33 := v.noise[x3+y3].Float64()

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
	xr := x.Remainder().Float64()
	yr := y.Remainder().Float64()
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

func (v *Value) CubicFloat(x, y float64) float64 {
	xf, yf := fixed.Float(x), fixed.Float(y)
	return v.Cubic(xf, yf)
}
