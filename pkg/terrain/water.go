package water

type Terrain struct {
	// Land is the height of solid ground above some arbitrary base value.
	Land []float64

	// Level is the height of the water above solid ground.
	Water []float64

	// VX is the velocity of water at any location along the X-axis.
	VX []float64

	// VY is the velocity of water at any location along the Y-axis.
	VY []float64
}
