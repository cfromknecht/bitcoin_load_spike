package bitcoin_load_spike

import "fmt"

type SpikeProfile struct {
	spikes []spike
}

func NewSpikeProfile(spikeMap map[float64]float64) *SpikeProfile {
	sp := new(SpikeProfile)
	sp.spikes = []spike{}

	for time, load := range spikeMap {
		sp.spikes = append(sp.spikes, spike{time, load})
	}

	return sp
}

type spike struct {
	time float64
	load float64
}

func (s spike) String() string {
	return fmt.Sprintf("%.4f:%.4f", s.load, s.time)
}

func (sp *SpikeProfile) valid() bool {
	previousTime := 0.0
	for i, spike := range sp.spikes {
		// First spike must be at time 0.0
		if i == 0 && spike.time != 0.0 {
			return false
		}
		// Check that all times and loads are valid
		if !validTime(spike.time) || !validLoad(spike.load) {
			return false
		}
		// Check that the times are in order
		if spike.time < previousTime {
			return false
		}
		previousTime = spike.time
	}
	return true
}

func (sp SpikeProfile) PrintProfile() {
	for _, spike := range sp.spikes {
		fmt.Println(fmt.Sprintf("    %3.f%%: %f", 100*spike.time, spike.load))
	}
}

func (sp SpikeProfile) CurrentLoad(percent float64) float64 {
	for i, spike := range sp.spikes {
		if percent == spike.time {
			return sp.spikes[i].load
		} else if percent < spike.time {
			return sp.spikes[i-1].load
		}
	}
	lastSpikeIndex := len(sp.spikes) - 1
	return sp.spikes[lastSpikeIndex].load
}

func (sp SpikeProfile) CurrentSpikeIndex(percent float64) int {
	for i, spike := range sp.spikes {
		if percent == spike.time {
			return i
		} else if percent < spike.time {
			return i - 1
		}
	}
	return len(sp.spikes) - 1
}

func validTime(t float64) bool {
	return t >= 0.0 && t <= 1.0
}

func validLoad(l float64) bool {
	return l >= 0.0
}
