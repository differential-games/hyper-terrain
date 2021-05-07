package fixed

type V4 [4]float64

func (v V4) TimesMatrix(m0, m1, m2, m3, m4, m5, m6, m7, m8, m9, m10, m11, m12, m13, m14, m15 float64) V4 {
	return V4{
		v[0]*m0 + v[1]*m4 + v[2]*m8 + v[3]*m12,
		v[0]*m1 + v[1]*m5 + v[2]*m9 + v[3]*m13,
		v[0]*m2 + v[1]*m6 + v[2]*m10 + v[3]*m14,
		v[0]*m3 + v[1]*m7 + v[2]*m11 + v[3]*m15,
	}
}

func (v V4) Dot(b0, b1, b2, b3 float64) float64 {
	return v[0]*b0 + v[1]*b1 + v[2]*b2 + v[3]*b3
}

// TimesM1 multiplies by the value matrix.
// See https://en.wikipedia.org/wiki/Bicubic_interpolation.
func TimesM1(v0, v1, v2, v3 float64) V4 {
	return V4{
		v0 - 3*v2 + 2*v3,
		3*v2 - 2*v3,
		v1 - 2*v2 + v3,
		-v2 + v3,
	}
}

// TimesM2 multiplies by the inverse of the value matrix.
// See https://en.wikipedia.org/wiki/Bicubic_interpolation.
func (v V4) TimesM2() V4 {
	return V4{
		v[0],
		v[2],
		-3*v[0] + 3*v[1] - 2*v[2] - v[3],
		2*v[0] - 2*v[1] + v[2] + v[3],
	}
}
