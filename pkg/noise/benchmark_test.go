package noise

import (
	"math/rand"
	"testing"
	"time"
)

func randP() pos {
	return p(rand.Float64()*float64(size), rand.Float64()*float64(size))
}

var px = randP()

func BenchmarkValue_Linear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n.Linear(px.x, px.y)
	}
}

func BenchmarkValue_Cubic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n.Cubic(px.x, px.y)
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
