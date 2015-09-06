package main

import (
	"flag"
	bls "github.com/cfromknecht/bitcoin_load_spike"
)

func parseFlags() (loadPercentage *float64, numBlocks, numSimulations *int64) {
	loadPercentage = flag.Float64("load", bls.DEFAULT_LOAD_PERCENTAGE, "percentage of maximum tps")
	numBlocks = flag.Int64("nb", bls.DEFAULT_NUM_BLOCKS, "number of blocks")
	numSimulations = flag.Int64("ns", bls.DEFAULT_NUM_SIMULATIONS, "number of simulations")

	flag.Parse()

	if loadPercentage == nil || numBlocks == nil || numSimulations == nil {
		panic("usage: [--load <load-percentage>] [--nb <num-blocks>] [--ns <num-simulations>]")
	}

	return
}

func main() {
	load, nb, ns := parseFlags()

	absoluteTps := *load * bls.BITCOIN_MAX_TPS
	sim := bls.NewLoadSpikeSimulation(absoluteTps, *nb, *ns)
	sim.Run()
}
