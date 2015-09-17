package main

import (
	"flag"
	bls "github.com/cfromknecht/bitcoin_load_spike"
	"runtime"
)

func parseFlags() (load, blockSize *float64, numBlocks, numIterations *int64) {
	load = flag.Float64("load", 0.0, "load percentage")
	blockSize = flag.Float64("bs", bls.DEFAULT_BLOCK_SIZE, "block size")
	numBlocks = flag.Int64("nb", bls.DEFAULT_NUM_BLOCKS, "number of blocks")
	numIterations = flag.Int64("ni", bls.DEFAULT_NUM_ITERATIONS, "number of iterations")

	flag.Parse()
	return
}

func main() {
	// Default to two processes, this will increase after further optimizations to the `createTxns` method
	runtime.GOMAXPROCS(2)

	load, bs, nb, ns := parseFlags()

	// Use constant `SpikeProfile` if `load` is set, otherwise use custom `SpikeProfile`
	var sp *bls.SpikeProfile
	if *load != 0.0 {
		sp = &bls.SpikeProfile{
			[]bls.Spike{
				bls.Spike{0.0, *load},
			},
		}
	} else {
		sp = &bls.SpikeProfile{
			[]bls.Spike{
				bls.Spike{0.0, 0.1},
				bls.Spike{0.33, 10.0},
				bls.Spike{0.67, 0.11},
			},
		}
	}

	// Run simulation with appropriate `SpikeProfile`
	bls.NewLoadSpikeSimulation(*bs, *nb, *ns).
		UseSpikeProfile(sp).
		AddCumulativeLogger("data/load-spike").
		//AddTimeSeriesLogger("data/load-spike").
		Run()
}
