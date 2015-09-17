# bitcoin_load_spike
Stochastic load spike modeling for bitcoin transactions

# Running
`go run run/main.go [--load <float>] [--bs <float>] [--nb <int>] [--ns <int>]`

`--load` the percentage of the maximum TPS of the bitcoin network. Setting this parameter will run the simulation with this single load for the entirety of the simulation.  Abscence of this parameter will cause the simulation to run the designation `SpikeProfile` described in `run/main.go`

`--bs` the maximum block size for the simulation

`--nb` number of blocks to create in a single iteration, roughly upperbounds the maximum time before transaction confirmation times are unrecorded by the simulation

`--ns` number of iterations to repeat using the above parameters, higher = more accurate

# Spike Profiles
Custom spike profiles can be defined in the `run/main.go` file.  Spikes are `(time, load)` pairs where `time` is the completion percentage [0,1) of the simulation (in terms of block numbers) and `load` is the percentage [0, infinity) of the maximum TPS (Transactions Per Second) of the network, currently set to 3.5 as per the Bitcoin Traffic Bulletin.

These spikes must occur in increasing order of time, the simulation will validate provided spike profiles and panic if they do not meet the above requirements requirements.

# Cumulative Logging
Data generated from each simulation is written to a file named `/data/load-spike-%f:%f-%d-%d.cl-dat`, where the format specifiers are replaced with the time, load, number of blocks, and number of iterations, respectively.  

Each files contains rows corresponding to `<bucket-number> | <log-of-txn-confirmation-time> | <probability> | <cumulative-probability>`.

# Plotting
`python plotter.py` will accrue all files in the `/data` folder with the format `load-spike-*.cl-dat` and attempt to plot them all in a single chart.  The resulting chart is then written to `/plots/load-spike-cumulatives.png`.

# Notes
The `time_series_logger.go` is still in the works, hence the lack of comments and support in the plotter.  Our goal is to first mimic the results seen in the Bitcoin Traffic Bulletin before introducing other improvements.
