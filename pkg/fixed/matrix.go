package fixed

// Matrix is an optimized 4x4 matrix.
type Matrix [16]float64

// Times is the inline matrix product of the two 4x4 matrices.
func (m Matrix) Times(b Matrix) Matrix {
	return Matrix{
		m[0]*b[0]+m[1]*b[4]+m[2]*b[8]+m[3]*b[12],
		m[0]*b[1]+m[1]*b[5]+m[2]*b[9]+m[3]*b[13],
		m[0]*b[2]+m[1]*b[6]+m[2]*b[10]+m[3]*b[14],
		m[0]*b[3]+m[1]*b[7]+m[2]*b[11]+m[3]*b[15],

		m[4]*b[0]+m[5]*b[4]+m[6]*b[8]+m[7]*b[12],
		m[4]*b[1]+m[5]*b[5]+m[6]*b[9]+m[7]*b[13],
		m[4]*b[2]+m[5]*b[6]+m[6]*b[10]+m[7]*b[14],
		m[4]*b[3]+m[5]*b[7]+m[6]*b[11]+m[7]*b[15],

		m[8]*b[0]+m[9]*b[4]+m[10]*b[8]+m[11]*b[12],
		m[8]*b[1]+m[9]*b[5]+m[10]*b[9]+m[11]*b[13],
		m[8]*b[2]+m[9]*b[6]+m[10]*b[10]+m[11]*b[14],
		m[8]*b[3]+m[9]*b[7]+m[10]*b[11]+m[11]*b[15],

		m[12]*b[0]+m[13]*b[4]+m[14]*b[8]+m[15]*b[12],
		m[12]*b[1]+m[13]*b[5]+m[14]*b[9]+m[15]*b[13],
		m[12]*b[2]+m[13]*b[6]+m[14]*b[10]+m[15]*b[14],
		m[12]*b[3]+m[13]*b[7]+m[14]*b[11]+m[15]*b[15],
	}
}

func (m Matrix) TimesV4(v V4) V4 {
	return V4{
		m[0]*v[0]+m[1]*v[1]+m[2]*v[2]+m[3]*v[3],
		m[4]*v[0]+m[5]*v[1]+m[6]*v[2]+m[7]*v[3],
		m[8]*v[0]+m[9]*v[1]+m[10]*v[2]+m[11]*v[3],
		m[12]*v[0]+m[13]*v[1]+m[14]*v[2]+m[15]*v[3],
	}
}

type V4 [4]float64

func (v V4) TimesMatrix(m Matrix) V4 {
	return V4{
		v[0]*m[0]+v[1]*m[4]+v[2]*m[8]+v[3]*m[12],
		v[0]*m[1]+v[1]*m[5]+v[2]*m[9]+v[3]*m[13],
		v[0]*m[2]+v[1]*m[6]+v[2]*m[10]+v[3]*m[14],
		v[0]*m[3]+v[1]*m[7]+v[2]*m[11]+v[3]*m[15],
	}
}

func (v V4) Dot(b V4) float64 {
	return v[0]*b[0] + v[1]*b[1] + v[2]*b[2] + v[3]*b[3]
}

// TimesM1 multiplies by the value matrix.
// See https://en.wikipedia.org/wiki/Bicubic_interpolation.
func (v V4) TimesM1() V4 {
	return V4{
		v[0]-3*v[2]+2*v[3],
		3*v[2]-2*v[3],
		v[1]-2*v[2]+v[3],
		-v[2]+v[3],
	}
}

// TimesM2 multiplies by the inverse of the value matrix.
// See https://en.wikipedia.org/wiki/Bicubic_interpolation.
func (v V4) TimesM2() V4 {
	return V4{
		v[0],
		v[2],
		-3*v[0]+3*v[1]-2*v[2]-v[3],
		2*v[0]-2*v[1]+v[2]+v[3],
	}
}
