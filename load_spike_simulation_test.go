package bitcoin_load_spike

import (
	"testing"
)

func TestNewLoadSpikeSimulation(t *testing.T) {
	expectedNumBlocks := int64(1000)
	expectedNumSimulations := int64(10000)

	sim := NewLoadSpikeSimulation(expectedNumBlocks, expectedNumSimulations)

	if sim.numBlocks != expectedNumBlocks {
		t.Error("Expected numBlocks to be", expectedNumBlocks, ", got", sim.numBlocks)
	}
	if sim.numIterations != expectedNumSimulations {
		t.Error("Expected numBlocks to be", expectedNumSimulations, ", got", sim.numIterations)
	}
	if sim.txnQ.headPtr != nil || sim.txnQ.tailPtr != nil {
		t.Error("Expected txnQ to be empty on initialization")
	}
	if sim.cacheQ.headPtr != nil || sim.cacheQ.tailPtr != nil {
		t.Error("Expected cacheQ to be empty on initialization")
	}
}
