package bitcoin_load_spike

// Assumptions
const BITCOIN_BLOCK_RATE float64 = 1.0 / 600.0              // 1 block every 10 minutes
const BITCOIN_TRANSACTION_SIZE int64 = (1024 * 1024) / 1200 // ~500 bytes
const BITCOIN_MAX_TPS float64 = 3.5                         // maximum number of txns per sec

// Default simulation parameters
const DEFAULT_LOAD_PERCENTAGE = 0.3    // Defaults to 30% of BITCOIN_MAX_TPS
const DEFAULT_NUM_BLOCKS = 1008        // One week of mining
const DEFAULT_NUM_SIMULATIONS = 100000 // Nice sample size

// Bucketing parameters for output
const NEGATIVE_ORDERS = 1
const POSITIVE_ORDERS = 10
const NUM_BUCKETS_PER_ORDER = 1000
const NUM_BUCKETS = (NUM_BUCKETS_PER_ORDER * (POSITIVE_ORDERS + NEGATIVE_ORDERS))

// Sampling parameters
const MAX_RAND = int64(2147483647)

// Error handling
func check(e error) {
	if e != nil {
		panic(e)
	}
}
