package bitcoin_load_spike

import "testing"

func TestNewLoadSpikeSimulation(t *testing.T) {
	expectedNumBlocks := int64(1000)
	expectedNumIterations := int64(10000)
	expectedBlockSize := DEFAULT_BLOCK_SIZE

	sim := NewLoadSpikeSimulation(expectedBlockSize, expectedNumBlocks, expectedNumIterations)

	if sim.blockSize != expectedBlockSize {
		t.Error("Expected blockSize to be", expectedBlockSize, ", got", sim.blockSize)
	}
	if sim.numBlocks != expectedNumBlocks {
		t.Error("Expected numBlocks to be", expectedNumBlocks, ", got", sim.numBlocks)
	}
	if sim.numIterations != expectedNumIterations {
		t.Error("Expected numIterations to be", expectedNumIterations, ", got", sim.numIterations)
	}
	if sim.spikeProfile != nil {
		t.Error("Expected spikeProfile to be nil, got", sim.spikeProfile)
	}
	if len(sim.loggers) != 0 {
		t.Error("Expected loggers to have length 0, got", len(sim.loggers))
	}
}

func TestUseSpikeProfile(t *testing.T) {
	expectedSpikes := []Spike{
		Spike{0.0, 0.1},
		Spike{0.5, 0.3},
	}
	expectedSpikeProfile := &SpikeProfile{
		Spikes: expectedSpikes,
	}

	sim := NewLoadSpikeSimulation(DEFAULT_BLOCK_SIZE, int64(1000), int64(1000)).
		UseSpikeProfile(expectedSpikeProfile)

	if sim.spikeProfile.Spikes[0].Percent != expectedSpikes[0].Percent || sim.spikeProfile.Spikes[0].Load != expectedSpikes[0].Load {
		t.Error("Expected spike", expectedSpikes[0], "got", sim.spikeProfile.Spikes[0])
	}
	if sim.spikeProfile.Spikes[1].Percent != expectedSpikes[1].Percent || sim.spikeProfile.Spikes[1].Load != expectedSpikes[1].Load {
		t.Error("Expected spike", expectedSpikes[1], "got", sim.spikeProfile.Spikes[1])
	}
}

func TestCreateTxns(t *testing.T) {
	expectedSpikeProfile := &SpikeProfile{
		[]Spike{
			Spike{0.0, 0.1},
			Spike{0.5, 0.3},
		},
	}

	expectedNumBlocks := int64(10)
	expectedNumIterations := int64(10)

	sim := NewLoadSpikeSimulation(DEFAULT_BLOCK_SIZE, expectedNumBlocks, expectedNumIterations).
		UseSpikeProfile(expectedSpikeProfile)

	pendingTxnChan := make(chan txn)
	readyChan := make(chan bool)
	blockNumChan := make(chan int64)

	go sim.createTxns(pendingTxnChan, readyChan, blockNumChan)

	for b := int64(0); b < expectedNumBlocks; b++ {
		blockNumChan <- b

		var tn txn
		if b == int64(0) {
			tn = <-pendingTxnChan
		} else {
			readyChan <- true
			tn = <-pendingTxnChan
		}

		// Two spikes, first half should be index 0, second half should be index 1
		if b < expectedNumBlocks/2 {
			if int64(tn.index) != 0 {
				t.Error("Expected spike index to be 0 got", tn.index)
			}
		} else {
			if int64(tn.index) != 1 {
				t.Error("Expected spike index to be 1 got", tn.index)
			}
		}

	}
	close(blockNumChan)
}

func TestCreateBlocks(t *testing.T) {
	expectedSpikeProfile := &SpikeProfile{
		[]Spike{
			Spike{0.0, 0.1},
			Spike{0.5, 0.3},
		},
	}

	expectedNumBlocks := int64(10)
	expectedNumIterations := int64(1)

	sim := NewLoadSpikeSimulation(DEFAULT_BLOCK_SIZE, expectedNumBlocks, expectedNumIterations).
		UseSpikeProfile(expectedSpikeProfile)

	pendingTxnChan := make(chan txn)
	readyChan := make(chan bool)
	blockNumChan := make(chan int64)

	go func(s *LoadSpikeSimulation, pChan chan txn, rChan chan bool, bChan chan int64) {
		currentTxnTimestamp := float64(0)
		numBlocks := int64(0)
		for {
			select {
			case i, ok := <-bChan:
				if !ok {
					// Check last iteration
					if numBlocks != expectedNumBlocks {
						t.Error("Expected number of blocks", expectedNumBlocks, "got", i+1)
					}

					close(pendingTxnChan)
					close(readyChan)
					return
				}
				if i == 0 {
					// Don't check on first iteration
					if numBlocks != 0 && numBlocks != expectedNumBlocks {
						t.Error("Expected number of blocks", expectedNumBlocks, "got", i+1)
					}
					numBlocks = 1

					currentTxnTimestamp = drawFromPoisson(0.35)
					pChan <- txn{currentTxnTimestamp, 0}
				} else {
					numBlocks++
				}

			case _ = <-readyChan:
				currentTxnTimestamp += drawFromPoisson(0.35)
				pChan <- txn{currentTxnTimestamp, 0}
			}
		}
	}(sim, pendingTxnChan, readyChan, blockNumChan)

	sim.createBlocks(pendingTxnChan, readyChan, blockNumChan)
}
