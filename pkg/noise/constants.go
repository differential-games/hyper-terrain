package noise

const (
	// shift is a compile-time constant representing the number of bits to shift 1 by to get
	// the size of noise to use. We want the noise size to be a power of 2 so that we can
	// use bit shifts instead of multiplication and division, and bitwise operations for
	// moduli.
	//
	// At 9, noise is 512x512. Higher numbers degrade performance, and all integers 1-9
	// have indistinguishable performance as the entire set of generated noise fits in a 2MiB L2
	// cache.
	shift uint8 = 9

	// revShift is a compile-time constant representing the number of bits to shift y from to
	// get it into the correct index position.
	revShift = 16 - shift

	// size is the scale of noise in units (for certain implementations) before it repeats.
	// We have it as a compile-time constant so types can use it as array lengths.
	size int = 1 << shift

	// intMask provides a convenient integer to take a bitwise-and with in order to perform
	// a cheap modulus with size.
	intMask = size - 1

	// size2 is the canonical square of size for use in 2D noise.
	size2 = size * size

	// int2Mask provides a convenient integer to take a bitwise-and with in order to perform
	// a cheap modulus with size2.
	int2Mask = size2 - 1 - intMask
)
