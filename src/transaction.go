package bitcoin_load_spike

const TRANSACTION_SIZE = 250 // 250 bytes

type txn struct {
	nextPtr *txn
	size int64
	time float64
}

type NewTxn(time float64) txn {
  return txn{nil, nil, TRANSACTION_SIZE, time}
}

type txnQueue struct {
	headPtr *txn
	tailPtr *txn
}

func NewTxnQueue() txnQueue {
	return txnQueue{nil, nil}
}

func (tq txnQueue) pushTxn(tPtr *txn) {
  if tq.headPtr == nil {
    tq.headPtr = tPtr
    tq.tailPtr = tPtr
  } else {
    tq.tailPtr.nextPtr = tPtr
    tq.tailPtr = tPtr
  }
}

func (tq txnQueue) popTxn() *txn {
  if tq.headPtr == nil {
    return nil
  } else {
    tPtr = tq.headPtr
    tq.headPtr = tq.headPtr.nextPtr

    // remove tail reference if list is empty
    if tq.headPtr == nil {
      tq.tailPtr = nil
    }
  }
}

