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
func newTxn(time float64) (txnPtr *txn) {
	txnPtr = new(txn)
	txnPtr.size = BITCOIN_TRANSACTION_SIZE
	txnPtr.time = time
	return
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
func (tq *txnQueue) pushTxn(txnPtr *txn) {
	txnPtr.nextPtr = nil
	if tq.headPtr == nil {
		tq.headPtr = txnPtr
		tq.tailPtr = txnPtr
	} else {
		tq.tailPtr.nextPtr = txnPtr
		tq.tailPtr = txnPtr
	}
}

/**
 * Removes a `txn` from the head of the `txnQueue`
 */
func (tq *txnQueue) popTxn() *txn {
	if tq.headPtr == nil {
		return nil
	} else {
		// Remove head element and advance references
		txnPtr := tq.headPtr
		tq.headPtr = tq.headPtr.nextPtr
		// Remove reference to current head
		txnPtr.nextPtr = nil

		// Remove tail reference if list is empty
		if tq.headPtr == nil {
			tq.tailPtr = nil
		}

		return txnPtr
	}
}

/**
 * Removes all `txn`s from the `txnQueue`
 */
func (tq *txnQueue) clear() {
	// Pop all elements from `txnQueue`
	currentTxnPtr := tq.popTxn()
	for currentTxnPtr != nil {
		currentTxnPtr = tq.popTxn()
	}
}
