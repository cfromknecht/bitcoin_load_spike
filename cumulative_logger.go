package bitcoin_load_spike

import (
	"fmt"
	"math"
)

/**
 * Stores the `cumulativePlot` for each spike in a `SpikeProfile` and the file
 * prefix for the output files.
 */
type CumulativeLogger struct {
	plots      []*cumulativePlot
	filePrefix string
}

/**
 * @return - The specified prefix for the output file.
 */
func (cl CumulativeLogger) FilePrefix() string {
	return cl.filePrefix
}

/**
 * @return - The file extension for `CumulativeLogger` output
 */
func (cl CumulativeLogger) FileExtension() string {
	return "cl-dat"
}

/**
 * Records the log of a `txn`s confirmation time into the appropriate bucket.
 *
 * @param blockTimestamp - The timestamp of the block that recorded `txn`
 * @param t - The `txn` that was recorded
 *
 */
func (cl *CumulativeLogger) Log(blockTimestamp float64, t txn) {
	// Caclulate the log of a txn's confirmation time
	age := blockTimestamp - t.time
	logAge := math.Log10(age)
	logAgeBucket := float64(NUM_BUCKETS_PER_ORDER) * logAge

	b := int64(math.Ceil(logAgeBucket))
	// Offset for negtive log values
	b += NEGATIVE_ORDERS * NUM_BUCKETS_PER_ORDER

	if b < 0 {
		b = 0
	}
	if b >= int64(len(cl.plots[t.index].buckets)) {
		// diff := b - int64(len(cl.plots[t.index].buckets) - 1)

		// bucketsExtension := make([]int64, diff)
		// cl.plots[t.index].buckets = append(cl.plots[t.index].buckets, bucketsExtension...)
		panic("Not enough buckets to record txn confirmation time.")
	}

	cl.plots[t.index].incrementBucket(b)
}

/**
 * Accumulates the file contents for all `cumulativePlot`s.
 *
 * @return - The file contents for each `Spike` in the `SpikeProfile`
 */
func (cl *CumulativeLogger) Outputs() (outputs []string) {
	for i, plot := range cl.plots {
		fmt.Println("[CumulativePlot]: generating cumulative plot data for spike", i)
		outputs = append(outputs, plot.output())
	}
	return
}

/**
 * Clears the logging state.
 */
func (cl *CumulativeLogger) Reset() {
	for i := range cl.plots {
		cl.plots[i] = newCumulativePlot()
	}
}

/**
 * Stores the buckets as an array of counters.  The number in each bucket
 * represents the number of txn's whose confirmation times fall within that bucket.
 * Also maintains a count of the total `txn`s recorded and the range of buckets
 * in use.
 */
type cumulativePlot struct {
	buckets        []int64
	smallestBucket int64
	largestBucket  int64
	txnCount       int64
}

/**
 * Initializes a new `cumulativePlot`
 *
 * @return - An empty `cumulativePlot`
 */
func newCumulativePlot() *cumulativePlot {
	return &cumulativePlot{
		buckets:        make([]int64, NUM_BUCKETS),
		smallestBucket: NUM_BUCKETS,
		largestBucket:  0,
		txnCount:       0,
	}
}

/**
 * Increments the bucket and total txn count. Also adjusts the range of used buckets.
 */
func (cp *cumulativePlot) incrementBucket(i int64) {
	cp.buckets[i]++
	cp.txnCount++

	if cp.largestBucket < i {
		cp.largestBucket = i
	}
	if cp.smallestBucket > i {
		cp.smallestBucket = i
	}
}

/**
 *  Returns a string representation of the plot to be written to a file.
 *
 * @return - The file contents for this spike's plot.
 */
func (cp *cumulativePlot) output() (fileContents string) {
	cumulativeTotal := float64(0.0)
	txnCountFloat := float64(cp.txnCount)

	for i, count := range cp.buckets[cp.smallestBucket:cp.largestBucket] {
		bucketCount := float64(count)
		cumulativeTotal += bucketCount

		fileContents += fmt.Sprintf("%d | %f | %f | %f\n",
			i,
			math.Pow(10.0, float64(i-(NEGATIVE_ORDERS*NUM_BUCKETS_PER_ORDER))/float64(NUM_BUCKETS_PER_ORDER)),
			bucketCount/txnCountFloat,
			cumulativeTotal/txnCountFloat)
	}
	return
}
