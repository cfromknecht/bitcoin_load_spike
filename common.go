package bitcoin_load_spike

// Assumptions
const BITCOIN_BLOCK_RATE = 1.0 / 600.0     // 1 block every 10 minutes
const BITCOIN_TRANSACTION_SIZE int64 = 500 // 500 bytes
const BITCOIN_MAX_TPS = 3.5                // maximum number of txns per sec

// Default simulation parameters
const DEFAULT_LOAD_PERCENTAGE = 0.3 // Defaults to 30% of BITCOIN_MAX_TPS
const DEFAULT_NUM_BLOCKS = 1000
const DEFAULT_NUM_SIMULATIONS = 100

// Bucketing parameters for output
const NEGATIVE_ORDERS = 1
const POSITIVE_ORDERS = 10
const NUM_BUCKETS_PER_ORDER = 1000
const NUM_BUCKETS = (NUM_BUCKETS_PER_ORDER * (POSITIVE_ORDERS + NEGATIVE_ORDERS))

// Sampling parameters
const MAX_SAMPLE_INT = int64(999999999)

// Error handling
func check(e error) {
	if e != nil {
		panic(e)
	}
}
