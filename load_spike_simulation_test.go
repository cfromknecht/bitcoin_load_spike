package bitcoin_load_spike

import (
	"testing"
)

func TestNewLoadSpikeSimulation(t *testing.T) {
	expectedTxnsPerSec := float64(3.5)
	expectedNumBlocks := int64(1000)
	expectedNumSimulations := int64(10000)

	sim := NewLoadSpikeSimulation(expectedTxnsPerSec, expectedNumBlocks, expectedNumSimulations)

	if sim.txnsPerSec != expectedTxnsPerSec {
		t.Error("Expected txnsPerSec to be", expectedTxnsPerSec, ", got", sim.txnsPerSec)
	}
	if sim.numBlocks != expectedNumBlocks {
		t.Error("Expected numBlocks to be", expectedNumBlocks, ", got", sim.numBlocks)
	}
	if sim.numSimulations != expectedNumSimulations {
		t.Error("Expected numBlocks to be", expectedNumSimulations, ", got", sim.numSimulations)
	}
	if len(sim.buckets) != NUM_BUCKETS {
		t.Error("Expected length of buckets to be", NUM_BUCKETS, ", got", len(sim.buckets))
	}
	if sim.smallestBucket != NUM_BUCKETS {
		t.Error("Expected smallestBucket to be", NUM_BUCKETS, ", got", sim.smallestBucket)
	}
	if sim.largestBucket != 0 {
		t.Error("Expected largestBucket to be", 0, ", got", sim.largestBucket)
	}
	if sim.txnCount != 0 {
		t.Error("Expected txnCount to be", 0, ", got", sim.txnCount)
	}
	if sim.txnQ.headPtr != nil || sim.txnQ.tailPtr != nil {
		t.Error("Expected txnQ to be empty on initialization")
	}
	if sim.cacheQ.headPtr != nil || sim.cacheQ.tailPtr != nil {
		t.Error("Expected cacheQ to be empty on initialization")
	}
}

func TestRecordAgeInBuckets(t *testing.T) {
	sim := NewLoadSpikeSimulation(1.0, 1000, 10000)

	expectedFirstBucket := int64(2699)
	sim.recordAgeInBuckets(50.0)

	// first addition should have same largest and smallest buckets
	if sim.smallestBucket != expectedFirstBucket {
		t.Error("Expected smallestBucket to be", expectedFirstBucket, ", got", sim.smallestBucket)
	}
	if sim.largestBucket != expectedFirstBucket {
		t.Error("Expected largestBucket to be", expectedFirstBucket, ", got", sim.largestBucket)
	}
	if sim.txnCount != 1 {
		t.Error("Expected txnCount to be 1, got", sim.txnCount)
	}

	expectedSecondBucket := int64(2779)
	sim.recordAgeInBuckets(60.0)

	// largest bucket should have increased, while smallest stays the same
	if sim.smallestBucket != expectedFirstBucket {
		t.Error("Expected smallestBucket to be", expectedFirstBucket, ", got", sim.smallestBucket)
	}
	if sim.largestBucket != expectedSecondBucket {
		t.Error("Expected largestBucket to be", expectedSecondBucket, ", got", sim.largestBucket)
	}
	if sim.txnCount != 2 {
		t.Error("Expected txnCount to be 2, got", sim.txnCount)
	}

	expectedThirdBucket := int64(2603)
	sim.recordAgeInBuckets(40.0)

	// smallest bucket should have decreased, while largest stays the same
	if sim.smallestBucket != expectedThirdBucket {
		t.Error("Expected smallestBucket to be", expectedThirdBucket, ", got", sim.smallestBucket)
	}
	if sim.largestBucket != expectedSecondBucket {
		t.Error("Expected largestBucket to be", expectedSecondBucket, ", got", sim.largestBucket)
	}
	if sim.txnCount != 3 {
		t.Error("Expected txnCount to be 3, got", sim.txnCount)
	}
}
