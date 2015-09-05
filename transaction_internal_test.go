package bitcoin_load_spike

import (
	"testing"
)

func TestNewTxn(t *testing.T) {
	expectedSize := TRANSACTION_SIZE
	expectedTime := 60.0

	txn := newTxn(expectedTime)

	// verify that field are set properly for new txn
	if txn.nextPtr != nil {
		t.Error("Expected nextPtr to be nil, got", txn.nextPtr)
	}
	if txn.size != expectedSize {
		t.Error("Expected size", expectedSize, "got", txn.size)
	}
	if txn.time != expectedTime {
		t.Error("Expected time", expectedTime, "got", txn.time)
	}
}

func TestNewTxnQueue(t *testing.T) {
	tq := newTxnQueue()

	// check that txnQueue is initialized with null references
	if tq.headPtr != nil {
		t.Error("Expected headPtr to be nil, got", tq.headPtr)
	}
	if tq.tailPtr != nil {
		t.Error("Expected tailPtr to be nil, got", tq.tailPtr)
	}
}

func TestPushTxn(t *testing.T) {
	expectedQueueLength := 10
	tq := newTxnQueue()

	// push a couple txns on to the queue
	for i := 1; i <= expectedQueueLength; i++ {
		txnPtr := newTxn(float64(i))
		tq.pushTxn(&txnPtr)
	}

	// iterate through list and count number of txns
	queueLength := 0
	currentPtr := tq.headPtr
	for i := 1; i <= expectedQueueLength; i++ {
		// check that txns are ordered properly
		if currentPtr.time != float64(i) {
			t.Error("Expected txn in queue to have time", float64(i), ", got", currentPtr.time)
		}
		// increment count and advance current pointer
		queueLength++
		currentPtr = currentPtr.nextPtr
	}

	// check that queue has the proper length
	if queueLength != expectedQueueLength {
		t.Error("Expected txnQueue length of", expectedQueueLength, ", got", queueLength)
	}
}

func TestPopTxn(t *testing.T) {
	expectedQueueLength := 10
	tq := newTxnQueue()

	// push a couple txns to the queue
	for i := 1; i <= expectedQueueLength; i++ {
		txnPtr := newTxn(float64(i))
		tq.pushTxn(&txnPtr)
	}

	// check that length of queue decreases with each pop
	for i := expectedQueueLength; i >= 1; i-- {
		txnPtr := tq.popTxn()
		// check that transactions are popped in reverse order
		expectedTime := float64(expectedQueueLength - i + 1)
		if txnPtr.time != expectedTime {
			t.Error("Expected popped txn to have time", expectedTime, ", got", txnPtr.time)
		}
	}

	// check that queue is empty again
	if tq.headPtr != nil {
		t.Error("Expected headPtr to be nil, got", tq.headPtr)
	}
	if tq.tailPtr != nil {
		t.Error("Expected tailPtr to be nil, got", tq.tailPtr)
	}
}