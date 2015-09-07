package bitcoin_load_spike

import (
	"math"
	"math/rand"
)

func drawFromPoisson(rate float64) float64 {
	r := rand.Float64()
	return float64(-math.Log(1.0-r) / rate)
}
