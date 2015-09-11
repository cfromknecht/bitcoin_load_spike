package bitcoin_load_spike

import "testing"

var validSpikeProfileTests = []struct {
	sp    SpikeProfile
	valid bool
	name  string
}{
	{
		SpikeProfile{
			[]Spike{Spike{0, .1}},
		},
		true,
		"initialization",
	},
	{
		SpikeProfile{
			[]Spike{Spike{0, .1}, Spike{.2, .1}},
		},
		true,
		"valid spike added",
	},
	{
		SpikeProfile{
			[]Spike{Spike{0, .1}, Spike{-.6, .1}},
		},
		false,
		"negative time",
	},
	{
		SpikeProfile{
			[]Spike{Spike{0, .1}, Spike{.3, -.1}},
		},
		false,
		"negative load",
	},
	{
		SpikeProfile{
			[]Spike{Spike{0, .1}, Spike{.5, .1}, Spike{.1, .1}},
		},
		false,
		"unordered spikes",
	},
}

func TestValidSpikeProfile(t *testing.T) {
	for _, test := range validSpikeProfileTests {
		// Expected outcome as a string
		expectedString := "valid"
		if test.valid == false {
			expectedString = "invalid"
		}

		if test.sp.valid() != test.valid {
			t.Error("Expected spike profile to be", expectedString, "for test", test.name)
		}
	}

}

func TestValidPercent(t *testing.T) {
	if !validPercent(.5) {
		t.Error("Expected .5 to be a valid time")
	}
	if !validPercent(0) {
		t.Error("Expected 0 to be a valid time")
	}
	if validPercent(1) {
		t.Error("Expected 1 to be an invalid time")
	}
	if validPercent(-.5) {
		t.Error("Expected -.5 to be an invalid time")
	}
	if validPercent(1.5) {
		t.Error("Expected 1.5 to be an invalid time")
	}
}

func TestValidLoad(t *testing.T) {
	if !validLoad(.5) {
		t.Error("Expected .5 to be a valid load")
	}
	if !validLoad(0) {
		t.Error("Expected 0 to be a valid load")
	}
	if validLoad(-.5) {
		t.Error("Expected -.5 to be a valid load")
	}
}

func TestCurrentLoad(t *testing.T) {
	sp := SpikeProfile{
		[]Spike{
			Spike{0.0, 0.1},
			Spike{0.1, 0.8},
			Spike{0.2, 0.2},
		},
	}

	if sp.currentLoad(0) != 0.1 {
		t.Error("Expected current load at 0 to be 0.1, got", sp.currentLoad(0))
	}
	if sp.currentLoad(0.05) != 0.1 {
		t.Error("Expected current load at 0.05 to be 0.1, got", sp.currentLoad(0.05))
	}
	if sp.currentLoad(0.1) != 0.8 {
		t.Error("Expected current load at 0.1 to be 0.8, got", sp.currentLoad(0.1))
	}
	if sp.currentLoad(0.15) != 0.8 {
		t.Error("Expected current load at 0.15 to be 0.8, got", sp.currentLoad(0.15))
	}
	if sp.currentLoad(0.2) != 0.2 {
		t.Error("Expected current load at 0.2 to be 0.2, got", sp.currentLoad(0.2))
	}
	if sp.currentLoad(0.5) != 0.2 {
		t.Error("Expected current load at 0.5 to be 0.2, got", sp.currentLoad(0.5))
	}
}

func TestCurrentIndex(t *testing.T) {
	sp := SpikeProfile{
		[]Spike{
			Spike{0.0, 0.1},
			Spike{0.1, 0.8},
			Spike{0.2, 0.2},
		},
	}

	if sp.currentSpikeIndex(0) != 0 {
		t.Error("Expected current load at 0 to be 0, got", sp.currentSpikeIndex(0))
	}
	if sp.currentSpikeIndex(0.05) != 0 {
		t.Error("Expected current load at 0.05 to be 0, got", sp.currentSpikeIndex(0.05))
	}
	if sp.currentSpikeIndex(0.1) != 1 {
		t.Error("Expected current load at 0.1 to be 1, got", sp.currentSpikeIndex(0.1))
	}
	if sp.currentSpikeIndex(0.15) != 1 {
		t.Error("Expected current load at 0.15 to be 1, got", sp.currentSpikeIndex(0.15))
	}
	if sp.currentSpikeIndex(0.2) != 2 {
		t.Error("Expected current load at 0.2 to be 2, got", sp.currentSpikeIndex(0.2))
	}
	if sp.currentSpikeIndex(0.5) != 2 {
		t.Error("Expected current load at 0.5 to be 2, got", sp.currentSpikeIndex(0.5))
	}
}
