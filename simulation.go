package bitcoin_load_spike

import (
	"fmt"
)

const BITCOIN_BLOCK_RATE = 1.0 / 600.0

var buckets = make([]int64, NUM_BUCKETS)
var smallestBucket = int64(NUM_BUCKETS)
var largestBucket = int64(0)
var numResults = int64(0)

func SimulateLoadSpikes(txnsPerSec float64, numBlocks, numSimulations int64) {
	divisor := numSimulations / 100
	if divisor == 0 {
		divisor = 1
	}

	for i = 0; i < numSimulations; i++ {
		simulateMining(txnsPerSec, numBlocks)

		if i%divisor == 0 {
			fmt.Println("[SimulateLoadSpike]:", i, "% complete")
		}
	}

	outputResults()
}

func simulateMining(txnsPerSec float64, numBlocks int64) {
	tq := newTxnQueue()

	firstTxnSecs := float64(0.0)
	cumulativeTime := float64(0.0)

	for i := int64(0); i < numBlocks; i++ {
		// time to mine the next block
		cumulativeTime += drawFromPoisson(BITCOIN_BLOCK_RATE)
		// create new transactions for this window
		firstTxnSecs = simulateTxns(firstTxnSecs, cumulativeTime, txnsPerSec, &tq)
		// consume transactions to be recorded in block
		createBlock(cumulativeTime, &tq)
	}
}

func simulateTxns(nextTxnSecs, miningEndTime, txnsPerSec float64, tqPtr *txnQueue) float64 {
	for {
		if miningEndTime < nextTxnSecs {
			return nextTxnSecs
		}

		txnPtr := newTxn(nextTxnSecs)
		tqPtr.pushTxn(txnPtr)

		nextTxnSecs += drawFromPoisson(txnsPerSec)
	}
}

func createBlock(blockTimestamp float64, tqPtr *txnQueue) (numTxns int64) {
	t := tqPtr.popTxn()
	if t == nil {
		return
	}

	remainingBlockSize := int64(1024 * 1024)
	for remainingBlockSize >= t.size {
		remainingBlockSize -= t.size
		numTxns++

		age := blockTimestamp - t.time
		logAge := math.Log10(age)
		// TODO(@cfromknecht) add to bucket

		if t.next == nil {
			return
		}

		t = tqPtr.popTxn()
	}
}

func outputResults() {

}
