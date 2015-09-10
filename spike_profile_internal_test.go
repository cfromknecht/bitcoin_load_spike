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

func TestValidTime(t *testing.T) {
	if !validTime(.5) {
		t.Error("Expected .5 to be a valid time")
	}
	if !validTime(0) {
		t.Error("Expected 0 to be a valid time")
	}
	if !validTime(1) {
		t.Error("Expected 1 to be a valid time")
	}
	if validTime(-.5) {
		t.Error("Expected -.5 to be a invalid time")
	}
	if validTime(1.5) {
		t.Error("Expected 1.5 to be a invalid time")
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

	if sp.CurrentLoad(0) != 0.1 {
		t.Error("Expected current load at 0 to be 0.1, got", sp.CurrentLoad(0))
	}
	if sp.CurrentLoad(0.05) != 0.1 {
		t.Error("Expected current load at 0.05 to be 0.1, got", sp.CurrentLoad(0.05))
	}
	if sp.CurrentLoad(0.1) != 0.8 {
		t.Error("Expected current load at 0.1 to be 0.8, got", sp.CurrentLoad(0.1))
	}
	if sp.CurrentLoad(0.15) != 0.8 {
		t.Error("Expected current load at 0.15 to be 0.8, got", sp.CurrentLoad(0.15))
	}
	if sp.CurrentLoad(0.2) != 0.2 {
		t.Error("Expected current load at 0.2 to be 0.2, got", sp.CurrentLoad(0.2))
	}
	if sp.CurrentLoad(0.5) != 0.2 {
		t.Error("Expected current load at 0.5 to be 0.2, got", sp.CurrentLoad(0.5))
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

	if sp.CurrentSpikeIndex(0) != 0 {
		t.Error("Expected current load at 0 to be 0, got", sp.CurrentSpikeIndex(0))
	}
	if sp.CurrentSpikeIndex(0.05) != 0 {
		t.Error("Expected current load at 0.05 to be 0, got", sp.CurrentSpikeIndex(0.05))
	}
	if sp.CurrentSpikeIndex(0.1) != 1 {
		t.Error("Expected current load at 0.1 to be 1, got", sp.CurrentSpikeIndex(0.1))
	}
	if sp.CurrentSpikeIndex(0.15) != 1 {
		t.Error("Expected current load at 0.15 to be 1, got", sp.CurrentSpikeIndex(0.15))
	}
	if sp.CurrentSpikeIndex(0.2) != 2 {
		t.Error("Expected current load at 0.2 to be 2, got", sp.CurrentSpikeIndex(0.2))
	}
	if sp.CurrentSpikeIndex(0.5) != 2 {
		t.Error("Expected current load at 0.5 to be 2, got", sp.CurrentSpikeIndex(0.5))
	}
}
