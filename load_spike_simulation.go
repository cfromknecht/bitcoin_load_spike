package bitcoin_load_spike

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

/**
 * `txn`
 *
 * Records time and spike index of the transaction's creation.
 */
type txn struct {
	time  float64
	index int
}

/**
 * `LoadSpikeSimulation`
 *
 * Stores simulation parameters and facilitiates logging during the simulation's
 * execution.
 */
type LoadSpikeSimulation struct {
	numBlocks     int64
	numIterations int64
	blockSize     float64
	spikeProfile  *SpikeProfile
	loggers       []Logger
}

/**
 * Initializes a new `LoadSpikeSimulation` with the simulation parameters.  By
 * default, the `spikeProfile` is not specified.
 *
 * @param bs - the maximum block size in bytes
 * @param nb - number of blocks to mine in a single iteration
 * @param ni - number of iterations to perform
 *
 * @return - the new `LoadSpikeSimulation`
 */
func NewLoadSpikeSimulation(bs float64, nb, ni int64) *LoadSpikeSimulation {
	return &LoadSpikeSimulation{
		numBlocks:     nb,
		numIterations: ni,
		blockSize:     bs,
		spikeProfile:  nil,
		loggers:       []Logger{},
	}
}

/**
 * Runs the simulation, printing the parameters and progress bar.  `Logger`s
 * accumulate data about the simulation and are printed after the simulation
 * terminates.  `Run` will panic if no `SpikeProfile` has been set.
 */
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

	// Reset loggers in case the simulation is reused
	for _, logger := range lss.loggers {
		logger.Reset()
	}
}

/**
 * Sets the simulations `spikeProfile`
 *
 * @param sp - The desired `SpikeProfile` for the simulation
 *
 * @return - The updated `LoadSpikeSimulation`
 */
func (lss *LoadSpikeSimulation) UseSpikeProfile(sp *SpikeProfile) *LoadSpikeSimulation {
	if sp == nil || !sp.valid() {
		panic("Cannot add invalid SpikeProfile to LoadSpikeSimulation")
	}
	// Add spike profile to simulation
	lss.spikeProfile = sp

	return lss
}

/**
 * Defines an interface for logging `txn`s and retrieving the outputs to be
 * written to files.
 */
type Logger interface {
	FilePrefix() string
	FileExtension() string
	Log( /* blockTimestamp */ float64, txn)
	Outputs() []string
	Reset()
}

/**
 * Adds a unique `TimeSeriesLogger` to the simulation's `loggers`
 *
 * @param prefix - The file prefix for writing the output file
 *
 * @return - The updated `LoadSpikeSimulation`
 */
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

/**
 * Adds a unique `CumulativeLogger` to the simulation's `loggers`
 *
 * @param prefix - The file prefix for writing the output file
 *
 * @return - The updated `LoadSpikeSimulation`
 */
func (lss *LoadSpikeSimulation) AddCumulativeLogger(prefix string) *LoadSpikeSimulation {
	if lss.spikeProfile == nil {
		panic("Cannot add CumulativeLogger without first setting a SpikeProfile")
	}

	// Create a plot record for each spike
	numPlots := len(lss.spikeProfile.Spikes)
	plots := make([]*cumulativePlot, numPlots)
	for i := range plots {
		plots[i] = newCumulativePlot()
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

/**
 * Spawns a routine to produce `txn`s, which are passed through a channel and
 * consumed on the main routine when they are added to a block.
 */
func (lss *LoadSpikeSimulation) simulateMining() {
	pendingTxnChan := make(chan txn)
	readyChan := make(chan bool)
	blockNumChan := make(chan int64)

	// Spawn routine to produce transactions
	go lss.createTxns(pendingTxnChan, readyChan, blockNumChan)
	// Consume transactions on main routine
	lss.createBlocks(pendingTxnChan, readyChan, blockNumChan)
}

/**
 * Produces transactions with timestamps drawn from a poisson distribution.  The
 * distribution is updated according to the `SpikeProfile`.  Transactions are
 * passed back through channels to be consumed in `createBlock`.
 *
 * @param pendingTxnChan - Channel for sending pending `txn`s to be consumed.
 * @param blockNumChan - Channel for receiving the current simultion's progress.
 *                       Used to determine the current load and spike index.
 */
func (lss *LoadSpikeSimulation) createTxns(pendingTxnChan chan txn, readyChan chan bool, blockNumChan chan int64) {
	var currentTPS float64
	var currentSpikeIndex int
	var currentTxnTimestamp float64

	for {
		select {
		case i, ok := <-blockNumChan:
			// Finished mining, simulation is complete
			if !ok {
				close(pendingTxnChan)
				close(readyChan)
				return
			}

			// Use percentage to get current tps (for poisson process) and spike
			// index (for logging)
			percent := float64(i) / float64(lss.numBlocks)
			// Percentage of BITCOIN_MAX_TPS
			currentLoad := lss.spikeProfile.currentLoad(percent)
			currentTPS = currentLoad * BITCOIN_MAX_TPS
			// Determines which log do eventually record the transaction under
			currentSpikeIndex = lss.spikeProfile.currentSpikeIndex(percent)

			// If starting new iteration, reset timestamp and send first txn
			if i == 0 {
				currentTxnTimestamp = drawFromPoisson(currentTPS)
				pendingTxnChan <- txn{currentTxnTimestamp, currentSpikeIndex}
			}
		case _ = <-readyChan:
			// Create and broadcast next txn
			currentTxnTimestamp += drawFromPoisson(currentTPS)
			pendingTxnChan <- txn{currentTxnTimestamp, currentSpikeIndex}
		}
	}
}

/**
 * Consumes `txn`s produced by `createTxn` and logs each one to the simulations
 * `loggers`.
 *
 * @param pendingTxnChan - Channel for receiving pending `txn`s to be consumed
 * @param blockNumChan - Channel for sending the current simultion's progress to
 *                       `createTxns`. Used to determine the current load and
 *                       spike index.
 */
func (lss *LoadSpikeSimulation) createBlocks(pendingTxnChan chan txn, readyChan chan bool, blockNumChan chan int64) {
	currentBlockTimestamp := float64(0)

	var t txn
	usePreviousTxn := false
	for i := int64(0); i < lss.numBlocks; i++ {
		blockNumChan <- i

		currentBlockTimestamp += drawFromPoisson(BITCOIN_BLOCK_RATE)
		remainingBlockSize := lss.blockSize

		// Must process `txn` before starting loop if last block was full
		if usePreviousTxn {
			remainingBlockSize -= BITCOIN_TRANSACTION_SIZE
			lss.logTxn(currentBlockTimestamp, t, readyChan)

			usePreviousTxn = false
		}

		for t = range pendingTxnChan {
			// If `txn` belongs in next block or doesn't fit in current block, process
			// later
			if t.time >= currentBlockTimestamp || remainingBlockSize < BITCOIN_TRANSACTION_SIZE {
				usePreviousTxn = true
				break
			}

			remainingBlockSize -= BITCOIN_TRANSACTION_SIZE
			lss.logTxn(currentBlockTimestamp, t, readyChan)
		}
	}

	// Terminates channels in createTxns
	close(blockNumChan)
}

/**
 * Logs a `txn` and the timestamp of the block in which it was recorded to the
 * simulations `logggers`
 *
 * @param blockTimestamp - Timestamp of the block that recorded `txn`
 * @param t - The `txn` that was consumed
 * @param readyChan - Channel for signaling when `createTxn` should send the next `txn`
 */
func (lss *LoadSpikeSimulation) logTxn(blockTimestamp float64, t txn, readyChan chan bool) {
	for _, logger := range lss.loggers {
		logger.Log(blockTimestamp, t)
	}
	readyChan <- true
}

/**
 * Obtains the outputs from each logger and writes them to their specified file.
 */
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

/**
 * Prints progress bar `|=========(10)=========(20)======...===|`
 */
func printProgessUpdate(i, divisor int64) {
	// Prints `[Progress]: `
	if i != 0 && i%(10*divisor) == 0 {
		fmt.Print(fmt.Sprintf("(%d)", i/divisor))
	} else if i%divisor == 0 {
		fmt.Print("=")
	}
}
