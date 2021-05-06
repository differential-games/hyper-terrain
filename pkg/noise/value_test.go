package noise

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var n = Value{}

func init() {
	n.Fill(rand.New(rand.NewSource(time.Now().UnixNano())))
}

type pos struct {
	x float64
	y float64
}

func (p pos) String() string {
	return fmt.Sprintf("(%.2f,%.2f)", p.x, p.y)
}

func p(x, y float64) pos {
	return pos{x: x, y: y}
}

func TestValue_V(t *testing.T) {
	tcs := []struct {
		name string
		p    pos
	}{
		{
			name: "works for 0, 0",
			p:    p(0, 0),
		},
		{
			name: "works for 0.5, 0.5",
			p:    p(0.5, 0.5),
		},
		{
			name: "works at noise boundary",
			p:    pos{x: float64(size), y: float64(size)},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			n.Linear(tc.p.x, tc.p.y)
		})
	}
}

func TestValue_V_Monotonic(t *testing.T) {
	tcs := []struct {
		name string
		p1   pos
		p2   pos
		p3   pos
	}{
		{
			name: "monotonic (0, 0) to (0, 1)",
			p1:   p(0, 0),
			p2:   p(0, rand.Float64()),
			p3:   p(0, 1),
		},
		{
			name: "monotonic (0, 1) to (1, 1)",
			p1:   p(0, 1),
			p2:   p(rand.Float64(), 1),
			p3:   p(1, 1),
		},
		{
			name: "monotonic (0, 0) to (1, 0)",
			p1:   p(0, 0),
			p2:   p(rand.Float64(), 0),
			p3:   p(1, 0),
		},
		{
			name: "monotonic (1, 0) to (1, 1)",
			p1:   p(1, 0),
			p2:   p(1, rand.Float64()),
			p3:   p(1, 1),
		},
		{
			name: "monotonic (0, SIZE-1) to (0, SIZE)",
			p1:   pos{x: 0, y: float64(size - 1)},
			p2:   pos{x: 0, y: float64(size-1) + rand.Float64()},
			p3:   pos{x: 0, y: float64(size)},
		},
		{
			name: "monotonic (0, SIZE) to (0, SIZE+1)",
			p1:   pos{x: (0), y: float64(size)},
			p2:   pos{x: (0), y: float64(size) + rand.Float64()},
			p3:   pos{x: (0), y: float64(size + 1)},
		},
		{
			name: "monotonic (SIZE, 0) to (SIZE+1, 0)",
			p1:   pos{x: float64(size), y: 0},
			p2:   pos{x: float64(size) + rand.Float64(), y: 0},
			p3:   pos{x: float64(size + 1), y: 0},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			v1 := n.Linear(tc.p1.x, tc.p1.y)
			v2 := n.Linear(tc.p2.x, tc.p2.y)
			v3 := n.Linear(tc.p3.x, tc.p3.y)

			if v1 < v3 {
				if v1 > v2 || v2 > v3 {
					t.Fatalf("behavior nonmonotonic from (%v:%f) to (%v:%f) to (%v:%f)",
						tc.p1, v1, tc.p2, v2, tc.p3, v3)
				}
			} else {
				if v1 < v2 || v2 < v3 {
					t.Fatalf("behavior nonmonotonic from (%v:%f) to (%v:%f) to (%v:%f)",
						tc.p1, v1, tc.p2, v2, tc.p3, v3)
				}
			}
		})
	}
}

func TestValue_V_Modulus(t *testing.T) {
	tcs := []struct {
		name string
		p1   pos
		p2   pos
	}{
		{
			name: "equivalent (0, 0) and (0, SIZE)",
			p1:   pos{x: 0, y: 0},
			p2:   pos{x: 0, y: float64(size)},
		},
		{
			name: "equivalent (0, 0) and (SIZE, 0)",
			p1:   pos{x: 0, y: 0},
			p2:   pos{x: float64(size), y: 0},
		},
		{
			name: "equivalent (0, 0) and (SIZE, SIZE)",
			p1:   pos{x: 0, y: 0},
			p2:   pos{x: float64(size), y: float64(size)},
		},
		{
			name: "equivalent (1, 1) and (1, SIZE+1)",
			p1:   pos{x: 0, y: 0},
			p2:   pos{x: 0, y: float64(size)},
		},
		{
			name: "equivalent (1, 1) and (SIZE+1, 1)",
			p1:   pos{x: 0, y: 0},
			p2:   pos{x: float64(size), y: 0},
		},
		{
			name: "equivalent (1, 1) and (SIZE+1, SIZE+1)",
			p1:   pos{x: 0, y: 0},
			p2:   pos{x: float64(size), y: float64(size)},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			v1 := n.Linear(tc.p1.x, tc.p1.y)
			v2 := n.Linear(tc.p2.x, tc.p2.y)

			if v1 != v2 {
				t.Fatalf("expected equivalent values (%v:%f) and (%v:%f)",
					tc.p1, v1, tc.p2, v2)
			}
		})
	}
}