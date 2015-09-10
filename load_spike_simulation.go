package bitcoin_load_spike

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

type LoadSpikeSimulation struct {
	// Simulation parameters
	numBlocks     int64
	numIterations int64
	blockSize     float64
	spikeProfile  *SpikeProfile
	// Logging state
	loggers []Logger
}

func NewLoadSpikeSimulation(bs float64, nb, ns int64) *LoadSpikeSimulation {
	return &LoadSpikeSimulation{
		numBlocks:     nb,
		numIterations: ns,
		blockSize:     bs,
		spikeProfile:  nil,
		loggers:       []Logger{},
	}
}

func (lss *LoadSpikeSimulation) Run() {
	if lss.spikeProfile == nil {
		panic("Cannot run LoadSpikeSimulation without a SpikeProfile")
	}

	// Make sure to seed our randomness
	rand.Seed(time.Now().UTC().UnixNano())

	// Print simulation parameters
	fmt.Println("[LoadSpikeSimulation]")
	fmt.Println("     iterations:", lss.numIterations)
	fmt.Println("     blocks/iteration:", lss.numBlocks)
	fmt.Println("     block size:", lss.blockSize)
	fmt.Println("[SpikeProfile]")
	lss.spikeProfile.PrintProfile()

	// Calculate divisor for progress bar
	divisor := lss.numIterations / 100
	if divisor == 0 {
		divisor = 1
	}

	// Run simulation
	fmt.Print("[Progress] |")
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
	pendingTxnChan := make(chan txn)
	cachedTxnChan := make(chan txn)
	blockNumChan := make(chan int64)

	// Spawn routine to produce transactions
	go lss.createTxns(pendingTxnChan, cachedTxnChan, blockNumChan)
	// Consume transactions on main routine
	lss.createBlocks(pendingTxnChan, cachedTxnChan, blockNumChan)
}

func (lss *LoadSpikeSimulation) createTxns(pendingTxnChan, cachedTxnChan chan txn, blockNumChan chan int64) {
	// createBlocks waits for channel to close before returning
	defer func() {
		close(pendingTxnChan)
	}()

	txnTimestamp := float64(0)

	// Simulation always starts at 0%
	currentSpikeIndex := lss.spikeProfile.CurrentSpikeIndex(0)
	currentLoad := lss.spikeProfile.CurrentLoad(0)
	currentTPS := currentLoad * BITCOIN_MAX_TPS

	for {
		select {
		case i := <-blockNumChan:
			// Starting new iteration
			if i == 0 {
				// Time till first txn
				txnTimestamp := drawFromPoisson(currentTPS)
				// Broadcast first txn
				pendingTxnChan <- newTxn(txnTimestamp, currentSpikeIndex)
			}

			// Use percentage to get current tps (for poisson) and spike index (for logging)
			percent := float64(i) / float64(lss.numBlocks)

			// Percentage of BITCOIN_MAX_TPS
			currentLoad := lss.spikeProfile.CurrentLoad(percent)
			currentTPS = currentLoad * BITCOIN_MAX_TPS

			// Determines which log do eventually record the transaction under
			currentSpikeIndex = lss.spikeProfile.CurrentSpikeIndex(percent)
		case txn, ok := <-cachedTxnChan:
			// Finished mining, simulation is complete
			if !ok {
				return
			}

			// Time till next txn
			txnTimestamp += drawFromPoisson(currentTPS)
			// Reuse transaction
			txn.time = txnTimestamp
			txn.index = currentSpikeIndex
			// Broadcast next txn
			pendingTxnChan <- txn
		}
	}
}

func (lss *LoadSpikeSimulation) createBlocks(pendingTxnChan, cachedTxnChan chan txn, blockNumChan chan int64) {
	blockTimestamp := float64(0)

	var t txn
	var usePreviousTxn = false
	for i := int64(0); i < lss.numBlocks; i++ {
		blockNumChan <- i

		blockTimestamp += drawFromPoisson(BITCOIN_BLOCK_RATE)
		remainingBlockSize := lss.blockSize

		// Must process `txn` before starting loop if last block was full
		if usePreviousTxn {
			remainingBlockSize -= BITCOIN_TRANSACTION_SIZE
			lss.logAndCacheTxn(txn, cachedTxnChan)

			usePreviousTxn = false
		}

		for t = range pendingTxnChan {
			// If `txn` belongs in next block or doesn't fit in current block, process
			// later
			if t.time >= blockTimestamp || remainingBlockSize < BITCOIN_TRANSACTION_SIZE {
				usePreviousTxn = true
				break
			}

			remainingBlockSize -= BITCOIN_TRANSACTION_SIZE
			lss.logAndCacheTxn(txn, cachedTxnChan)
		}
	}

	// Terminates channels and waits for other routine to close the pendingTxnChan
	// before returning
	close(blockNumChan)
	close(cachedTxnChan)
	<-pendingTxnChan
}

func (lss *LoadSpikeSimulation) logAndCacheTxn(t txn, cachedTxnChan chan txn) {
	// Log the txn
	for _, logger := range lss.loggers {
		logger.Log(blockTimestamp, t.time, t.index)
	}
	// Return txn to other routine
	cachedTxnChan <- t
}

func (lss *LoadSpikeSimulation) UseSpikeProfile(sp *SpikeProfile) *LoadSpikeSimulation {
	if sp == nil || !sp.valid() {
		panic("Cannot add invalid SpikeProfile to LoadSpikeSimulation")
	}
	// Add spike profile to simulation
	lss.spikeProfile = sp

	return lss
}

func (lss *LoadSpikeSimulation) AddTimeSeriesLogger(prefix string) *LoadSpikeSimulation {
	// Create new time series logger
	tsLogger := &TimeSeriesLogger{
		plot:          newTimeSeriesPlot(),
		secsPerBucket: 60.0,
		filePrefix:    prefix,
	}

	// Append logger to loggers
	lss.loggers = append(lss.loggers, tsLogger)

	return lss
}

func (lss *LoadSpikeSimulation) AddCumulativeLogger(prefix string) *LoadSpikeSimulation {
	if lss.spikeProfile == nil {
		panic("Cannot add CumulativeLogger without first setting a SpikeProfile")
	}

	// Create a plot record for each spike
	numPlots := len(lss.spikeProfile.Spikes)
	plots := []*cumulativePlot{}
	for i := 0; i < numPlots; i++ {
		plots = append(plots, newCumulativePlot())
	}

	// Build logger
	cLogger := &CumulativeLogger{
		plots:      plots,
		filePrefix: prefix,
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
			filename += "-" + lss.spikeProfile.Spikes[i].String()
			filename += fmt.Sprintf("-%d-%d", lss.numBlocks, lss.numIterations)
			filename += "." + logger.FileExtension()
			// Write file contents to filename
			err := ioutil.WriteFile(filename, []byte(fileContents), 0644)
			check(err)
		}
	}
}
