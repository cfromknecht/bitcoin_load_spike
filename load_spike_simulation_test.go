package bitcoin_load_spike

import (
	"testing"
)

func TestNewLoadSpikeSimulation(t *testing.T) {
	expectedNumBlocks := int64(1000)
	expectedNumSimulations := int64(10000)
	expectedBlockSize := DEFAULT_BLOCK_SIZE

	sim := NewLoadSpikeSimulation(expectedBlockSize, expectedNumBlocks, expectedNumSimulations)

	if sim.blockSize != expectedBlockSize {
		t.Error("Expected blockSize to be", expectedBlockSize, ", got", sim.blockSize)
	}
	if sim.numBlocks != expectedNumBlocks {
		t.Error("Expected numBlocks to be", expectedNumBlocks, ", got", sim.numBlocks)
	}
	if sim.numIterations != expectedNumSimulations {
		t.Error("Expected numBlocks to be", expectedNumSimulations, ", got", sim.numIterations)
	}
}
