package bitcoin_load_spike

// Assumptions
const BITCOIN_BLOCK_RATE = 1.0 / 600.0     // 1 block every 10 minutes
const BITCOIN_TRANSACTION_SIZE int64 = 250 // 250 bytes

// Default simulation parameters
const DEFAULT_TXNS_PER_SEC = 1.0
const DEFAULT_NUM_BLOCKS = 1000
const DEFAULT_NUM_SIMULATIONS = 100

// Bucketing parameters for output
const NEGATIVE_ORDERS = 1
const POSITIVE_ORDERS = 10
const NUM_BUCKETS_PER_ORDER = 1000
const NUM_BUCKETS = (NUM_BUCKETS_PER_ORDER * (POSITIVE_ORDERS + NEGATIVE_ORDERS))

// Error handling
func check(e error) {
	if e != nil {
		panic(e)
	}
}
