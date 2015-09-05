package bitcoin_load_spike

/**
 * Record for storing a linked list of transactions
 */
type txn struct {
	nextPtr *txn
	size    int64 // static for now, added to support dynamic sizing in the future
	time    float64
}

/**
 * Creates a new `txn` at a specified `time`
 */
func newTxn(time float64) txn {
	return txn{nil, BITCOIN_TRANSACTION_SIZE, time}
}

/**
 * Queue of `txn`s, stored as a linked list
 */
type txnQueue struct {
	headPtr *txn
	tailPtr *txn
}

/**
 * Creates an empty queue of `txn`s
 */
func newTxnQueue() txnQueue {
	return txnQueue{nil, nil}
}

/**
 * Adds a `txn` to the tail of the `txnQueue`
 */
func (tq *txnQueue) pushTxn(tPtr *txn) {
	if tq.headPtr == nil {
		tq.headPtr = tPtr
		tq.tailPtr = tPtr
	} else {
		tq.tailPtr.nextPtr = tPtr
		tq.tailPtr = tPtr
	}
}

/**
 * Removes a `txn` from the head of the `txnQueue`
 */
func (tq *txnQueue) popTxn() *txn {
	if tq.headPtr == nil {
		return nil
	} else {
		tPtr := tq.headPtr
		tq.headPtr = tq.headPtr.nextPtr

		// remove tail reference if list is empty
		if tq.headPtr == nil {
			tq.tailPtr = nil
		}

		return tPtr
	}
}

/**
 * Removes all `txn`s from the `txnQueue`
 */
func (tq *txnQueue) clear() {
	// remove references between current transactions
	currentTxnPtr := tq.headPtr
	for currentTxnPtr != nil {
		nextPtr := currentTxnPtr.nextPtr
		currentTxnPtr.nextPtr = nil
		currentTxnPtr = nextPtr
	}
	// remove references to head and tail
	tq.headPtr = nil
	tq.tailPtr = nil
}
