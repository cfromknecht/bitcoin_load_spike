package bitcoin_load_spike

import (
	"fmt"
	_ "math"
)

type timeSeriesPlot struct {
	buckets        []float64
	counts         []int64
	smallestBucket int64
	largestBucket  int64
}

func newTimeSeriesPlot() *timeSeriesPlot {
	return &timeSeriesPlot{
		buckets:        make([]float64, NUM_BUCKETS),
		counts:         make([]int64, NUM_BUCKETS),
		smallestBucket: NUM_BUCKETS,
		largestBucket:  0,
	}
}

func (tsp *timeSeriesPlot) updateBucket(i int64, age float64) {
	total := tsp.buckets[i]*float64(tsp.counts[i]) + age
	tsp.counts[i]++
	tsp.buckets[i] = total / float64(tsp.counts[i])

	if tsp.largestBucket < i {
		tsp.largestBucket = i
	}
	if tsp.smallestBucket > i {
		tsp.smallestBucket = i
	}
}

func (tsp *timeSeriesPlot) output() (fileContents string) {
	for i, avgTxnTime := range tsp.buckets[0 : NUM_BUCKETS-1] {
		fileContents += fmt.Sprintf("%d | %f\n", i, avgTxnTime)
	}
	return
}

type TimeSeriesLogger struct {
	plot          *timeSeriesPlot
	secsPerBucket float64
	filePrefix    string
}

func (tsl TimeSeriesLogger) FilePrefix() string {
	return tsl.filePrefix
}

func (tsl TimeSeriesLogger) FileExtension() string {
	return "tsl-dat"
}

func (tsl *TimeSeriesLogger) Log(blockTimestamp, txnTimestamp float64, spikeNumber int) {
	age := blockTimestamp - txnTimestamp

	b := int64(txnTimestamp / tsl.secsPerBucket)
	// Extend buckets and counts if necessary
	if b >= int64(len(tsl.plot.buckets)) {
		diff := b - int64(len(tsl.plot.buckets)-1)
		// Extend buckets
		bucketsExtension := make([]float64, diff)
		tsl.plot.buckets = append(tsl.plot.buckets, bucketsExtension...)
		// Extend counts
		countsExtension := make([]int64, diff)
		tsl.plot.counts = append(tsl.plot.counts, countsExtension...)
	}

	tsl.plot.updateBucket(b, age)
}

func (tsl *TimeSeriesLogger) Outputs() (outputs []string) {
	fmt.Println("[TimeSeriesLogger]: generating time series plot")
	outputs = append(outputs, tsl.plot.output())
	return
}

func (tsl *TimeSeriesLogger) Reset() {
	tsl.plot = newTimeSeriesPlot()
}
