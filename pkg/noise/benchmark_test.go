package noise

import (
	"math/rand"
	"testing"
	"time"
)

func randP() (float64, float64) {
	return rnd.Float64()*float64(size), rnd.Float64()*float64(size)
}

func Benchmark_RandP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randP()
	}
}

func BenchmarkValue_Linear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n.Linear(randP())
	}
}

func BenchmarkValue_Cubic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n.Cubic(randP())
	}
}

var (
	n2  = Value{}
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func BenchmarkValue_Fill(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n2.Fill(rnd)
	}
}
