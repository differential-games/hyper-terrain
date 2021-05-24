package noise

import (
	"math/rand"
	"testing"
	"time"
)

func randP(rnd *rand.Rand) (x, y float64) {
	return rnd.Float64() * float64(size), rnd.Float64() * float64(size)
}

func Benchmark_RandP(b *testing.B) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.ResetTimer()

	// Benchmark computing random points for tests.
	for i := 0; i < b.N; i++ {
		randP(rnd)
	}
}

func BenchmarkValue_Linear(b *testing.B) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.ResetTimer()

	n := Value{}
	n.Fill(rand.New(rand.NewSource(time.Now().UnixNano())))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n.Linear(randP(rnd))
	}
}

func BenchmarkValue_Cubic(b *testing.B) {
	n := Value{}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	n.Fill(rnd)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n.Cubic(randP(rnd))
	}
}

func BenchmarkValue_Fill(b *testing.B) {
	n2 := Value{}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n2.Fill(rnd)
	}
}
