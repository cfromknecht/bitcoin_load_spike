package main

import (
	"flag"
	bls "github.com/cfromknecht/bitcoin_load_spike"
)

func parseFlags() (blockSize *float64, numBlocks, numIterations *int64) {
	blockSize = flag.Float64("bs", bls.DEFAULT_BLOCK_SIZE, "block size")
	numBlocks = flag.Int64("nb", bls.DEFAULT_NUM_BLOCKS, "number of blocks")
	numIterations = flag.Int64("ni", bls.DEFAULT_NUM_ITERATIONS, "number of iterations")

	flag.Parse()

	if blockSize == nil || numBlocks == nil || numIterations == nil {
		panic("usage: [--bs <block-size>] [--nb <num-blocks>] [--ni <num-iterations>]")
	}

	return
}

func main() {
	bs, nb, ns := parseFlags()

	sp := &bls.SpikeProfile{
		Spikes: []bls.Spike{
			bls.Spike{0.0, 1.0},
		},
	}

	bls.NewLoadSpikeSimulation(*bs, *nb, *ns).
		UseSpikeProfile(sp).
		AddCumulativeLogger("data/load-spike").
		AddTimeSeriesLogger("data/load-spike").
		Run()
}
