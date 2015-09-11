package bitcoin_load_spike

import (
	"math"
	"math/rand"
)

/**
 * Returns a random sample from a poisson distribution for a given `rate`.
 *
 * @return - A poisson sample for the given `rate`
 */
func drawFromPoisson(rate float64) float64 {
	r := rand.Float64()
	return -float64(math.Log(1.0-r) / rate)
}
