package bitcoin_load_spike

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

type LoadSpikeSimulation struct {
	// simulation parameters
	numBlocks     int64
	numIterations int64
	spikeProfile  *SpikeProfile
	// simulation state
	txnQ    txnQueue
	cacheQ  txnQueue
	loggers []Logger
}

func NewLoadSpikeSimulation(nb, ns int64) *LoadSpikeSimulation {
	return &LoadSpikeSimulation{
		numBlocks:     nb,
		numIterations: ns,
		spikeProfile:  nil,
		txnQ:          newTxnQueue(),
		cacheQ:        newTxnQueue(),
		loggers:       []Logger{},
	}
}

func (lss *LoadSpikeSimulation) Run() {
	if lss.spikeProfile == nil {
		panic("Cannot run LoadSpikeSimulation without a SpikeProfile")
	}

	// Print parameters
	fmt.Println("[LoadSpikeSimulation]: simulating", lss.numIterations, "iterations with", lss.numBlocks, "blocks each")
	fmt.Println("[SpikeProfile]:")
	lss.spikeProfile.PrintProfile()

	// Make sure we have a strong seed
	rand.Seed(time.Now().UTC().UnixNano())

	divisor := lss.numIterations / 100
	if divisor == 0 {
		divisor = 1
	}

	// Run simulation
	fmt.Print("[Progress]: |")
	for i := int64(0); i < lss.numIterations; i++ {
		lss.simulateMining()
		printProgessUpdate(i, divisor)
	}
	fmt.Println("|")

	lss.OutputResults()

	for _, logger := range lss.loggers {
		logger.Reset()
	}
}

func printProgessUpdate(i, divisor int64) {
	// Prints `[Progress]: |=========(10)=========(20)======`
	if i != 0 && i%(10*divisor) == 0 {
		fmt.Print(fmt.Sprintf("(%d)", i/divisor))
	} else if i%divisor == 0 {
		fmt.Print("=")
	}
}

func (lss *LoadSpikeSimulation) simulateMining() {
	firstTxnSecs := float64(0.0)
	cumulativeTime := float64(0.0)

	for i := int64(0); i < lss.numBlocks; i++ {
		percent := float64(i) / float64(lss.numBlocks)
		// time to mine the next block
		cumulativeTime += drawFromPoisson(BITCOIN_BLOCK_RATE)
		// create new transactions for this window
		firstTxnSecs = lss.simulateTxns(firstTxnSecs, cumulativeTime, percent)
		// consume as many transactions as possible into the next block
		lss.createBlock(cumulativeTime, percent)
	}

	// Move all remaining txns to the `cacheQ`
	currentTxnPtr := lss.txnQ.popTxn()
	for currentTxnPtr != nil {
		lss.cacheQ.pushTxn(currentTxnPtr)
		currentTxnPtr = lss.txnQ.popTxn()
	}
}

func (lss *LoadSpikeSimulation) simulateTxns(nextTxnSecs, miningEndTime float64, percent float64) float64 {
	for {
		if miningEndTime < nextTxnSecs {
			return nextTxnSecs
		}

		// Try to utilize previously allocated txn
		txnPtr := lss.cacheQ.popTxn()
		if txnPtr != nil {
			txnPtr.time = nextTxnSecs
		} else {
			// Otherwise create new txn
			txnPtr = newTxn(nextTxnSecs)
		}
		lss.txnQ.pushTxn(txnPtr)

		// Calculate TPS from current percentage complete
		currentLoad := lss.spikeProfile.CurrentLoad(percent)
		currentTPS := currentLoad * BITCOIN_MAX_TPS

		nextTxnSecs += drawFromPoisson(currentTPS)
	}
}

func (lss *LoadSpikeSimulation) createBlock(blockTimestamp float64, percent float64) {
	txnPtr := lss.txnQ.popTxn()
	if txnPtr == nil {
		return
	}

	remainingBlockSize := int64(1024 * 1024)
	for remainingBlockSize >= txnPtr.size {
		remainingBlockSize -= txnPtr.size

		// Log the results
		for _, logger := range lss.loggers {
			currentSpikeIndex := lss.spikeProfile.CurrentSpikeIndex(percent)
			logger.Log(blockTimestamp, txnPtr.time, currentSpikeIndex)
		}

		// Cache this transaction for later
		lss.cacheQ.pushTxn(txnPtr)

		// Pop and return if queue is empty
		txnPtr = lss.txnQ.popTxn()
		if txnPtr == nil {
			return
		}
	}
}

func (lss *LoadSpikeSimulation) UseSpikeProfile(sp *SpikeProfile) *LoadSpikeSimulation {
	if sp == nil || !sp.valid() {
		panic("Cannot add invalid SpikeProfile to LoadSpikeSimulation")
	}
	// Add spike profile to simulation
	lss.spikeProfile = sp

	return lss
}

func (lss *LoadSpikeSimulation) AddCumulativeLogger(prefix string) *LoadSpikeSimulation {
	if lss.spikeProfile == nil {
		panic("Cannot add CumulativeLogger without first setting a SpikeProfile")
	}

	// Create a plot record for each spike
	numPlots := len(lss.spikeProfile.spikes)
	plots := []*cumulativePlot{}
	for i := 0; i < numPlots; i++ {
		plots = append(plots, newCumulativePlot())
	}

	// Build logger
	cLogger := &CumulativeLogger{
		plots,
		prefix,
	}

	// Append logger to loggers
	lss.loggers = append(lss.loggers, cLogger)

	return lss
}

func (lss *LoadSpikeSimulation) OutputResults() {
	// Create output for each logger
	for _, logger := range lss.loggers {
		// Create file prefix to dump results
		filePrefix := logger.FilePrefix()

		// Get each file contents and write to file
		for i, fileContents := range logger.Outputs() {
			// Create full filename
			filename := filePrefix
			filename += "-" + lss.spikeProfile.spikes[i].String()
			filename += fmt.Sprintf("-%d-%d", lss.numBlocks, lss.numIterations)
			filename += "." + logger.FileExtension()
			// Write file contents to filename
			err := ioutil.WriteFile(filename, []byte(fileContents), 0644)
			check(err)
		}
	}
}
