package bitcoin_load_spike

import (
	"fmt"
	"math"
)

type cumulativePlot struct {
	buckets        []int64
	smallestBucket int64
	largestBucket  int64
	txnCount       int64
}

func newCumulativePlot() *cumulativePlot {
	return &cumulativePlot{
		buckets:        make([]int64, NUM_BUCKETS),
		smallestBucket: NUM_BUCKETS,
		largestBucket:  0,
		txnCount:       0,
	}
}

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

type CumulativeLogger struct {
	plots      []*cumulativePlot
	filePrefix string
}

func (cl CumulativeLogger) FilePrefix() string {
	return cl.filePrefix
}

func (cl CumulativeLogger) FileExtension() string {
	return "cl-dat"
}

func (cl *CumulativeLogger) Log(blockTimestamp, txnTimestamp float64, spikeNumber int) {
	age := blockTimestamp - txnTimestamp
	logAge := math.Log10(age)
	logAgeBucket := float64(NUM_BUCKETS_PER_ORDER) * logAge

	b := int64(math.Ceil(logAgeBucket))
	b += NEGATIVE_ORDERS * NUM_BUCKETS_PER_ORDER

	if b < 0 {
		b = 0
	}
	if b >= NUM_BUCKETS {
		panic("Not enough buckets to record txn confirmation time.")
	}

	cl.plots[spikeNumber].incrementBucket(b)
}

func (cl *CumulativeLogger) Outputs() (outputs []string) {
	for i, plot := range cl.plots {
		fmt.Println("[CumulativePlot]: generating cumulative plot data for spike", i)
		outputs = append(outputs, plot.output())
	}
	return
}

func (cl *CumulativeLogger) Reset() {
	for i := range cl.plots {
		cl.plots[i] = newCumulativePlot()
	}
}
