package bitcoin_load_spike

import (
	"crypto/rand"
	"math"
	"math/big"
)

func drawFromPoisson(rate float64) float64 {
	r := randomPercentage()
	return -float64(math.Log(1.0-r) / rate)
}

func randomPercentage() float64 {
	// set max bound
	var maxInt int64 = 999999999
	max := *big.NewInt(maxInt)

	// generate random number
	rBig, err := rand.Int(rand.Reader, &max)
	check(err)

	// calculate percentage
	return float64(rBig.Uint64()) / (float64(maxInt) + 1)
}
