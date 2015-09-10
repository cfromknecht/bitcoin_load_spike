package bitcoin_load_spike

import "testing"

func TestFilePrefix(t *testing.T) {
	expectedPrefix := "prefix"

	cl := CumulativeLogger{
		[]*cumulativePlot{},
		expectedPrefix,
	}

	if cl.FilePrefix() != expectedPrefix {
		t.Error("Expected file prefix", expectedPrefix, ", got", cl.FilePrefix())
	}
}

func TestFileExtension(t *testing.T) {
	expectedExtension := "cl-dat"

	cl := CumulativeLogger{
		[]*cumulativePlot{},
		"",
	}

	if cl.FileExtension() != expectedExtension {
		t.Error("Expected file extension", expectedExtension, ", got", cl.FileExtension())
	}
}

var logTests = []struct {
	blockTimestamp float64
	txnTimestamp   float64
	expectedBucket int64
	shouldRecover  bool
}{
	{
		0.0,
		0.0,
		0,
		false,
	},
	{
		10.0,
		0.0,
		2000,
		false,
	},
	{
		10000.0,
		0,
		5000,
		false,
	},
	{
		100000000000000000000, // Some very high number
		0,
		0, // Not used, test for panicking instead
		true,
	},
}

func TestLog(t *testing.T) {
	for _, test := range logTests {
		// Test for panicking if index would be out of bounds
		if test.shouldRecover {
			defer func() {
				if r := recover(); r != nil {
					if r != "Not enough buckets to record txn confirmation time." {
						t.Error("Panic message different from expected, got", r)
					}
				} else {
					t.Error("Expected log to panic if bucket index would cause array out of bounds error")
				}
			}()
		}

		cl := CumulativeLogger{
			[]*cumulativePlot{newCumulativePlot()},
			"",
		}
		cl.Log(test.blockTimestamp, test.txnTimestamp, 0)

		if cl.plots[0].buckets[test.expectedBucket] != 1 {
			// find actual bucket
			var actualBucket int
			for i, count := range cl.plots[0].buckets {
				if count == 1 {
					actualBucket = i
					break
				}
			}
			t.Error("Expected bucket", test.expectedBucket, "to be incremented for block timestamp", test.blockTimestamp, "and txn timestamp", test.txnTimestamp, ", got", actualBucket)
		}
	}
}

func TestOutput(t *testing.T) {
	expectedOutput := "0 | 0.100000 | 0.400000 | 0.400000\n"
	cl := CumulativeLogger{
		[]*cumulativePlot{newCumulativePlot()},
		"",
	}
	for i := float64(0); i < 5; i++ {
		cl.Log(1000.0, i, 0)
	}

	output := cl.Outputs()[0]

	if output != expectedOutput {
		t.Error("Expected output '", expectedOutput, "', got '", output, "'")
	}
}
