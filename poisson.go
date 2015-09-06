package bitcoin_load_spike

import (
	"crypto/rand"
	"math"
	"math/big"
)

func drawFromPoisson(rate float64) float64 {
	r := randomPercentage()
	return float64(-math.Log(1.0-r) / rate)
}

func randomPercentage() float64 {
	bigMaxIntPtr := big.NewInt(MAX_SAMPLE_INT)

	// generate random number
	rBig, err := rand.Int(rand.Reader, bigMaxIntPtr)
	check(err)

	// calculate percentage
	return float64(rBig.Uint64()) / (float64(MAX_SAMPLE_INT) + 1.0)
}
