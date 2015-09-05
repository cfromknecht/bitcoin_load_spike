#! /usr/bin/env python

import sys
import re
import argparse
from matplotlib import pyplot



def main(args = sys.argv[1:]):
    """
    Plot results from either Bitcoin Traffic Bulletin or this repository.
    """
    parse_args(args)

    times, freqs, cumulatives = parse_stream(sys.stdin)
    make_plot(times, freqs, cumulatives)


def parse_args(args):
    p = argparse.ArgumentParser(description=main.__doc__)
    return p.parse_args(args)


def parse_stream(f):
    linergx = re.compile(r'^[\d]+ \| (\d+\.\d+) \| (\d+\.\d+) \| (\d+\.\d+)$')

    times, freqs, cumulatives = [], [], []

    for ix, line in enumerate(f):
        m = linergx.match(line)
        try:
            assert m is not None
            [t, f, c] = [ float(g) for g in m.groups() ]
        except Exception as e:
            sys.stderr.write(
                'Could not parse line {}: {}\n  Input: {!r}\n'
                .format(ix+1, e, line))
        else:
            times.append(t)
            freqs.append(f)
            cumulatives.append(c)

    return times, freqs, cumulatives


def make_plot(times, freqs, cumulatives):
    pyplot.plot(times, freqs)
    plotresult='./plot.png'
    print('Writing plot: {!r}'.format(plotresult))
    pyplot.savefig(plotresult)


if __name__ == '__main__':
    main()
