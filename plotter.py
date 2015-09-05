#! /usr/bin/env python

import sys
import os
import re
import glob
import argparse
from matplotlib import pyplot



def main(args = sys.argv[1:]):
    """
    Plot results from either Bitcoin Traffic Bulletin or this repository.
    """
    opts = parse_args(args)

    generate_spike_plots(opts.datadir, opts.plotdir)


def parse_args(args):
    reporoot = os.path.abspath(os.path.join(sys.argv[0], '..'))

    class defaults:
        datadir = os.path.join(reporoot, 'data')
        plotdir = os.path.join(reporoot, 'plots')

    p = argparse.ArgumentParser(description=main.__doc__)

    p.add_argument('--data-dir',
                   dest='datadir',
                   default=defaults.datadir,
                   help='Data directory. Default: {!r}'.format(defaults.datadir))

    p.add_argument('--plot-dir',
                   dest='plotdir',
                   default=defaults.plotdir,
                   help='Data directory. Default: {!r}'.format(defaults.plotdir))

    return p.parse_args(args)


def generate_spike_plots(datadir, plotdir):
    times, datasets = parse_spike_data(datadir)
    plot_spike_data(plotdir, times, datasets)


def parse_spike_data(datadir):
    fnamergx = re.compile(r'^.*/load-spike-(\d+\.\d+)-(\d+)-(\d+).dat$')

    fixedparams = None
    datasets = {}

    for path in glob.glob('{}/load-spike-*.dat'.format(datadir)):
        m = fnamergx.match(path)
        try:
            rate, blocks, sims = m.groups()
            rate = float(rate)
            blocks = int(blocks)
            sims = int(sims)
        except Exception as e:
            warnuser('Could not parse parameters from path {!r}: {}', path, e)
            continue

        with file(path, 'r') as f:
            times, freqs, cumulatives = parse_stream(f)

        if fixedparams is None:
            # The only independent variable should be rate, and only
            # freqs/cumulatives should vary dependently:
            fixedparams = (blocks, sims, times)
        else:
            params = (blocks, sims, times)
            if fixedparams != params:
                warnuser(
                    'Skipping plot {!r} which has params {!r}, expected {!r}.',
                    path, params, fixedparams)
                continue

        datasets[rate] = (freqs, cumulatives)

    return (times, datasets)


def plot_spike_data(plotdir, times, datasets):
    pyplot.figure()

    for rate, (_, cumulatives) in sorted(datasets.items()):
        pyplot.plot(times, cumulatives, label='{}'.format(rate))

    plotpath = os.path.join(plotdir, 'load-spike-cumulatives.png')
    print('Writing plot: {!r}'.format(plotpath))
    pyplot.savefig(plotpath)


def parse_stream(f):
    linergx = re.compile(r'^[\d]+ \| (\d+\.\d+) \| (\d+\.\d+) \| (\d+\.\d+)$')

    times, freqs, cumulatives = [], [], []

    for ix, line in enumerate(f):
        m = linergx.match(line)
        try:
            assert m is not None
            [t, f, c] = [ float(g) for g in m.groups() ]
        except Exception as e:
            warnuser('Could not parse line {}: {}\n  Input: {!r}\n', ix+1, e, line)
        else:
            times.append(t)
            freqs.append(f)
            cumulatives.append(c)

    return times, freqs, cumulatives


def warnuser(tmpl, *args):
    sys.stderr.write(tmpl.format(*args) + '\n')


if __name__ == '__main__':
    main()
