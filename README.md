# bitcoin_load_spike
Stochastic load spike modeling for bitcoin transactions

# Running
`go run run/main.go [--load <float>] [--nb <int>] [--ns <int>]`

`--load` percentage of the Bitcoin's maximum TPS (Transactions Per Second), which is currently ~ 3.5tps

`--nb` number of blocks to create in a single iteration, roughly upperbounds the maximum time before transaction confirmation times are unrecorded by the simulation

`--ns` number of iterations to repeat using the above parameters, higher = more accurate

# Output
Data generated from each simulation is written to a file named `/data/load-spike-%f-%d-%d.dat`, where the format specifiers are replaced with the load, number of blocks, and number of iterations, respectively.  

Each files contains rows corresponding to `<bucket-number> | <log-of-txn-confirmation-time> | <probability> | <cumulative-probability>`.

# Plotting
`python plotter.py` will accrue all files in the `/data` folder with the format `load-spike-*.dat` and attempt to plot them all in a single chart.  The resulting chart is then written to `/plots/load-spike-cumulatives.png`.
