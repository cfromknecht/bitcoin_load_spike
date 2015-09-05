package bitcoin_load_spike

import (
	"fmt"
	"math"
)

const DEFAULT_TPS = 1.0
const DEFAULT_NB = 1000
const DEFAULT_NS = 100

const NEGATIVE_ORDERS = 1
const POSITIVE_ORDERS = 10
const NUM_BUCKETS_PER_ORDER = 1000
const NUM_BUCKETS = (NUM_BUCKETS_PER_ORDER * (POSITIVE_ORDERS + NEGATIVE_ORDERS))

const BITCOIN_BLOCK_RATE = 1.0 / 600.0 // 1 block every 10 minutes

var buckets = make([]int64, NUM_BUCKETS)
var smallestBucket = int64(NUM_BUCKETS)
var largestBucket = int64(0)
var numResults = int64(0)

func SimulateLoadSpikes(txnsPerSec float64, numBlocks, numSimulations int64) {
	divisor := numSimulations / 100
	if divisor == 0 {
		divisor = 1
	}

	for i := int64(0); i < numSimulations; i++ {
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
		tqPtr.pushTxn(&txnPtr)

		nextTxnSecs += drawFromPoisson(txnsPerSec)
	}
}

func createBlock(blockTimestamp float64, tqPtr *txnQueue) (numTxns int64) {
	txnPtr := tqPtr.popTxn()
	if txnPtr == nil {
		return
	}

	remainingBlockSize := int64(1024 * 1024)
	for remainingBlockSize >= txnPtr.size {
		remainingBlockSize -= txnPtr.size
		numTxns++

		age := blockTimestamp - txnPtr.time
		math.Log10(age)
		// TODO(@cfromknecht) add to bucket

		if txnPtr.nextPtr == nil {
			return
		}

		txnPtr = tqPtr.popTxn()
	}

	return
}

func outputResults() {

}
