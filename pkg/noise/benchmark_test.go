package noise

import (
	"math/rand"
	"testing"
	"time"
)

func randP() pos {
	return p(rand.Float64()*float64(size), rand.Float64()*float64(size))
}

func Benchmark_RandP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randP()
	}
}

func BenchmarkValue_Linear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pt := randP()
		n.Linear(pt.x, pt.y)
	}
}

func BenchmarkValue_Cubic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pt := randP()
		n.Cubic(pt.x, pt.y)
	}
}

var (
	n2  = Value{}
	src = rand.NewSource(time.Now().UnixNano())
)

func BenchmarkValue_Fill(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n2.Fill(src)
	}
}
