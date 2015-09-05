package bitcoin_load_spike_test

import (
	bls "github.com/cfromknecht/bitcoin_load_spike"
	"testing"
)

func TestNewTxn(t *testing.T) {
	expectedSize := bls.TRANSACTION_SIZE
	expectedTime := 60.0

	t := bls.NewTxn(expectedTime)

	if t.nextPtr != nil {
		t.Error("Expected nextPtr to be nil, got %p", t.nextPtr)
	}

	if t.size != expectedSize {
		t.Error("Expected size %d, got %d", expectedSize, t.size)
	}

	if t.time != expectedTime {
		t.Error("Expected time %f, got %f", expectedTime, t.time)
	}
}
