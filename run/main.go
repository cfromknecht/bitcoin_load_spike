package main

import (
	"flag"
	bls "github.com/cfromknecht/bitcoin_load_spike"
)

func parseFlags() (numBlocks, numIterations *int64) {
	numBlocks = flag.Int64("nb", bls.DEFAULT_NUM_BLOCKS, "number of blocks")
	numIterations = flag.Int64("ni", bls.DEFAULT_NUM_ITERATIONS, "number of iterations")

	flag.Parse()

	if numBlocks == nil || numIterations == nil {
		panic("usage: [--nb <num-blocks>] [--ni <num-iterations>]")
	}

	return
}

func main() {
	nb, ns := parseFlags()

	sp := bls.NewSpikeProfile(map[float64]float64{
		0.0: 0.1,
		0.2: 10,
		0.4: 0.1,
	})

	bls.NewLoadSpikeSimulation(*nb, *ns).
		UseSpikeProfile(sp).
		AddCumulativeLogger("data/load-spike").
		Run()
}
