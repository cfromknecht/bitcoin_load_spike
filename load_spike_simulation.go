package bitcoin_load_spike

import (
	"fmt"
	"io/ioutil"
	"math"
)

type LoadSpikeSimulation struct {
	// simulation parameters
	txnsPerSec     float64
	numBlocks      int64
	numSimulations int64
	// simulation state
	buckets        []int64
	smallestBucket int64
	largestBucket  int64
	txnCount       int64
	txnQ           txnQueue
}

func NewLoadSpikeSimulation(tps float64, nb, ns int64) *LoadSpikeSimulation {
	return &LoadSpikeSimulation{
		txnsPerSec:     tps,
		numBlocks:      nb,
		numSimulations: ns,
		buckets:        make([]int64, NUM_BUCKETS),
		smallestBucket: NUM_BUCKETS,
		largestBucket:  0,
		txnCount:       0,
		txnQ:           newTxnQueue(),
	}
}

func (lss *LoadSpikeSimulation) Run() {
	fmt.Println("[LoadSpikeSimulation] starting simulation with tps:",
		lss.txnsPerSec, "nb:", lss.numBlocks, "ns:", lss.numSimulations)

	divisor := lss.numSimulations / 100
	if divisor == 0 {
		divisor = 1
	}

	for i := int64(0); i < lss.numSimulations; i++ {
		lss.simulateMining()

		if i%divisor == 0 {
			fmt.Println("[LoadSpikeSimulation]:", i, "% complete")
		}
	}

	lss.outputResults()
}

func (lss *LoadSpikeSimulation) simulateMining() {
	firstTxnSecs := float64(0.0)
	cumulativeTime := float64(0.0)

	for i := int64(0); i < lss.numBlocks; i++ {
		// time to mine the next block
		cumulativeTime += drawFromPoisson(BITCOIN_BLOCK_RATE)
		// create new transactions for this window
		firstTxnSecs = lss.simulateTxns(firstTxnSecs, cumulativeTime)
		// consume as many transactions as possible into the next block
		lss.createBlock(cumulativeTime)
	}
}

func (lss *LoadSpikeSimulation) simulateTxns(nextTxnSecs, miningEndTime float64) float64 {
	for {
		if miningEndTime < nextTxnSecs {
			return nextTxnSecs
		}

		txnPtr := newTxn(nextTxnSecs)
		lss.txnQ.pushTxn(&txnPtr)

		nextTxnSecs += drawFromPoisson(lss.txnsPerSec)
	}
}

func (lss *LoadSpikeSimulation) createBlock(blockTimestamp float64) {
	txnPtr := lss.txnQ.popTxn()
	if txnPtr == nil {
		return
	}

	remainingBlockSize := int64(1024 * 1024)
	for remainingBlockSize >= txnPtr.size {
		remainingBlockSize -= txnPtr.size

		// time from transaction creation to being recorded in this block
		age := blockTimestamp - txnPtr.time
		lss.recordAgeInBuckets(age)

		if txnPtr.nextPtr == nil {
			return
		}

		txnPtr = lss.txnQ.popTxn()
	}
}

func (lss *LoadSpikeSimulation) recordAgeInBuckets(age float64) {
	logAge := math.Log10(age)
	logAgeBucket := float64(NUM_BUCKETS_PER_ORDER) * logAge

	b := int64(math.Ceil(logAgeBucket))
	b += NEGATIVE_ORDERS * NUM_BUCKETS_PER_ORDER

	if b < 0 {
		b = 0
	}

	// increment bucket for age and total txn count
	lss.buckets[b]++
	lss.txnCount++

	// update used bucket range
	if lss.largestBucket < b {
		lss.largestBucket = b
	}
	if lss.smallestBucket > b {
		lss.smallestBucket = b
	}
}

func (lss *LoadSpikeSimulation) outputResults() {
	loadPercentage := lss.txnsPerSec / BITCOIN_MAX_TPS

	filename := fmt.Sprintf("data/load-spike-%f-%d-%d.dat", loadPercentage, lss.numBlocks, lss.numSimulations)
	fileContents := ""

	cumulativeTotal := float64(0.0)
	txnCountFloat := float64(lss.txnCount)
	for i, count := range lss.buckets[lss.smallestBucket:lss.largestBucket] {
		bucketCount := float64(count)
		cumulativeTotal += bucketCount

		fileContents += fmt.Sprintf("%d | %f | %f | %f\n",
			i,
			math.Pow(10.0, float64(i-(NEGATIVE_ORDERS*NUM_BUCKETS_PER_ORDER))/float64(NUM_BUCKETS_PER_ORDER)),
			bucketCount/txnCountFloat,
			cumulativeTotal/txnCountFloat)
	}

	err := ioutil.WriteFile(filename, []byte(fileContents), 0644)
	check(err)
}
