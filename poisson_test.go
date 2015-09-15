package bitcoin_load_spike

import (
	"fmt"
	"testing"
)

var poissonTests = []float64{
	0.01,
	0.5,
	1.0,
	10.0,
}

func TestPoisson(t *testing.T) {
	for _, load := range poissonTests {
		total := float64(0)
		for i := 0; i < 100000; i++ {
			total += drawFromPoisson(float64(load))
		}
		fmt.Println("Total average for rate:", load, ", got", total/100000.0)
	}
}
