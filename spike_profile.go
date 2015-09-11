package bitcoin_load_spike

import "fmt"

/**
 * `SpikeProfile`
 *
 * Defines the load at any percentage of the simulation's completion by
 * interpolating the most recent spike for a given percentage.
 *
 * Example:
 * The `SpikeProfile` below can be defined by using 3 `Spike`s. One at 0%, 20%,
 * and 40%.  The loads are then interpretted until the beginning of the next
 * spike or the end of the simulation.
 *
 *      ___
 *     |   |
 *  ___|   |_________
 * |   |   |
 * 0% 20% 40%       100%
 */
type SpikeProfile struct {
	Spikes []Spike
}

/**
 * `Spike`
 *
 * Defines the load at a given percentage
 */
type Spike struct {
	Percent float64
	Load    float64
}

/**
 * Returns a string representation of a `Spike`
 *
 * @return - "<percent>:<load>"
 */
func (s Spike) String() string {
	return fmt.Sprintf("%.4f:%.4f", s.Percent, s.Load)
}

/**
 * Iterates over a `SpikeProfile` and prints each `Spike`s string representation
 */
func (sp SpikeProfile) PrintProfile() {
	for _, spike := range sp.Spikes {
		fmt.Println(fmt.Sprintf("    %3.f%%: %f", 100*spike.Percent, spike.Load))
	}
}

/**
 * Caclulates the current load given the percentage complete.
 *
 * @param percent - The simulation's completion percentage.
 *
 * @return - The network load at `percent`
 */
func (sp SpikeProfile) currentLoad(percent float64) float64 {
	currentSpikeIndex := sp.currentSpikeIndex(percent)
	return sp.Spikes[currentSpikeIndex].Load
}

/**
 * Caclulates the current spike index given the simulation's completion percentage.
 *
 * @param percent - The simulation's completion percentage.
 *
 * @return - The index of the current spike.
 */
func (sp SpikeProfile) currentSpikeIndex(percent float64) int {
	for i, spike := range sp.Spikes {
		if percent == spike.Percent {
			// Percentages match, return this index
			return i
		} else if percent < spike.Percent {
			// Simulation has not reached this spike yet, return pervious index
			return i - 1
		}
	}
	// Otherwise return the last spike's index
	return len(sp.Spikes) - 1
}

/**
 * Verifies that a `SpikeProfile` is valid for use in the simulation.  Checks
 * that all percentages are in [0, 1) and tht all loads are greater than 0.
 * Also checks that the percentages are ordered properly.
 */
func (sp *SpikeProfile) valid() bool {
	previousTime := 0.0
	for i, spike := range sp.Spikes {
		// First spike must be at time 0.0
		if i == 0 && spike.Percent != 0 {
			return false
		}
		// Check that all times and loads are valid
		if !validPercent(spike.Percent) || !validLoad(spike.Load) {
			return false
		}
		// Check that the times are in order
		if spike.Percent < previousTime {
			return false
		}
		previousTime = spike.Percent
	}
	return true
}

/**
 * Checks that `t` is in the range [0, 1)
 *
 * @return - Whether `t` is a valid percent
 */
func validPercent(t float64) bool {
	return t >= 0.0 && t < 1.0
}

/**
 * Checks that `l` is greater than 0
 *
 * @ return - Whether `t` is valid load
 */
func validLoad(l float64) bool {
	return l >= 0.0
}
