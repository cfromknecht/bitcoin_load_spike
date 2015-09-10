package bitcoin_load_spike

import "fmt"

type SpikeProfile struct {
	Spikes []Spike
}

type Spike struct {
	Time float64
	Load float64
}

func (s Spike) String() string {
	return fmt.Sprintf("%.4f:%.4f", s.Load, s.Time)
}

func (sp *SpikeProfile) valid() bool {
	previousTime := 0.0
	for i, spike := range sp.Spikes {
		// First spike must be at time 0.0
		if i == 0 && spike.Time != 0 {
			return false
		}
		// Check that all times and loads are valid
		if !validTime(spike.Time) || !validLoad(spike.Load) {
			return false
		}
		// Check that the times are in order
		if spike.Time < previousTime {
			return false
		}
		previousTime = spike.Time
	}
	return true
}

func (sp SpikeProfile) PrintProfile() {
	for _, spike := range sp.Spikes {
		fmt.Println(fmt.Sprintf("    %3.f%%: %f", 100*spike.Time, spike.Load))
	}
}

func (sp SpikeProfile) CurrentLoad(percent float64) float64 {
	for i, spike := range sp.Spikes {
		if percent == spike.Time {
			return sp.Spikes[i].Load
		} else if percent < spike.Time {
			return sp.Spikes[i-1].Load
		}
	}
	lastSpikeIndex := len(sp.Spikes) - 1
	return sp.Spikes[lastSpikeIndex].Load
}

func (sp SpikeProfile) CurrentSpikeIndex(percent float64) int {
	for i, spike := range sp.Spikes {
		if percent == spike.Time {
			return i
		} else if percent < spike.Time {
			return i - 1
		}
	}
	return len(sp.Spikes) - 1
}

func validTime(t float64) bool {
	return t >= 0.0 && t <= 1.0
}

func validLoad(l float64) bool {
	return l >= 0.0
}
