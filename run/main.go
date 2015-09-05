package main

import (
	"flag"
	bls "github.com/cfromknecht/bitcoin_load_spike"
)

func parseFlags() (txnsPerSec *float64, numBlocks, numSimulations *int64) {
	txnsPerSec = flag.Float64("tps", bls.DEFAULT_TPS, "transactions per second")
	numBlocks = flag.Int64("nb", bls.DEFAULT_NB, "number of blocks")
	numSimulations = flag.Int64("ns", bls.DEFAULT_NS, "number of simulations")

	flag.Parse()

	if txnsPerSec == nil || numBlocks == nil || numSimulations == nil {
		panic("usage: bitcoin-load-spike <transactions-per-second> <num-blocks> <num-simulations>")
	}

	return
}

func main() {
	tps, nb, ns := parseFlags()

	sim := bls.NewLoadSpikeSimulation(*tps, *nb, *ns)
	sim.Run()
}
