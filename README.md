# gowc
[![CI](https://github.com/svaloumas/gowc/actions/workflows/ci.yml/badge.svg)](https://github.com/svaloumas/gowc/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/svaloumas/gowc/branch/main/graph/badge.svg?token=7DvuWdQZPr)](https://codecov.io/gh/svaloumas/gowc)
[![Go Report Card](https://goreportcard.com/badge/github.com/svaloumas/gowc)](https://goreportcard.com/report/github.com/svaloumas/gowc)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/svaloumas/gowc/blob/main/LICENSE)

Just another GNU [`wc`](https://www.gnu.org/software/coreutils/manual/html_node/wc-invocation.html#wc-invocation) clone, written in Go.

## Overview

`gowc` is a simple, zero-dependency command line tool for counting bytes, characters, words and newlines in each given file.
It leverages the language's built-in support for concurrency by processing the given input files in chunks. The buffer size of each chunk is configurable
and can be set via `-bs, --buffer-size` flag. It reads one chunk ahead while processing the previously read one.

## Installation

```
make build-linux
```

Other available options: `build-mac`, `build-win`

## Usage

By default, `gowc` will count lines, words, and bytes. You can specify the counters you'd like by using the available flags and options from the table below.

| Flag | Description |
| ---- | ----------- |
| -c, --bytes | Print the byte counts |
| -m, --chars | Print the character counts |
| -l, --lines | Print the newline counts | 
| -l, --lines | Print the newline counts | 
| -w, --words | Print the word counts | 
| -L, --max-line-length | Print the length of the longest line | 
| -h, --help | Display help and exit | 
| -V, --version | Output version information and exit | 

| Option | Description |
| ------ | ----------- |
| -bs, --buffer-size   | Configure the buffer size of each chunk to be processed (defaults to 4096) | 
| --files-from <file>  | Read input from the files specified by a newline-terminated list of filenames in the given file | 

```shell
gowc [FLAGS] [OPTIONS] [FILE]...
```

## Performance

[`hyperfine`](https://github.com/sharkdp/hyperfine) is used to perform the benchmarks. The file used is a 595MB CSV with 5m rows.

```bash
# New lines only count
$ hyperfine  --warmup 3 './gowc -l -bs 1000000 ./5mSalesRecords.csv' 'wc -l ./5mSalesRecords.csv'                                                                        [±main ●]
Benchmark 1: ./gowc -l -bs 1000000 ./5mSalesRecords.csv
  Time (mean ± σ):     160.2 ms ±   6.5 ms    [User: 118.5 ms, System: 126.6 ms]
  Range (min … max):   148.8 ms … 167.4 ms    17 runs
 
Benchmark 2: wc -l ./5mSalesRecords.csv
  Time (mean ± σ):     494.3 ms ±  12.3 ms    [User: 397.0 ms, System: 93.8 ms]
  Range (min … max):   480.8 ms … 517.6 ms    10 runs
 
Summary
  './gowc -l -bs 1000000 ./5mSalesRecords.csv' ran
    3.08 ± 0.15 times faster than 'wc -l ./5mSalesRecords.csv'

# Default lines, words and bytes count
hyperfine  --warmup 3 './gowc -bs 1000000 ./5mSalesRecords.csv' 'wc ./5mSalesRecords.csv'                                                                              [±main ●]
Benchmark 1: ./gowc -bs 1000000 ./5mSalesRecords.csv
  Time (mean ± σ):      1.542 s ±  0.008 s    [User: 1.554 s, System: 0.239 s]
  Range (min … max):    1.532 s …  1.559 s    10 runs
 
Benchmark 2: wc ./5mSalesRecords.csv
  Time (mean ± σ):      2.045 s ±  0.009 s    [User: 1.946 s, System: 0.097 s]
  Range (min … max):    2.033 s …  2.058 s    10 runs
 
Summary
  './gowc -bs 1000000 ./5mSalesRecords.csv' ran
    1.33 ± 0.01 times faster than 'wc ./5mSalesRecords.csv'

# Word only count
$ hyperfine  --warmup 3 './gowc -w -bs 1000000 ./5mSalesRecords.csv' 'wc -w ./5mSalesRecords.csv'                                                                        [±main ●]
Benchmark 1: ./gowc -w -bs 1000000 ./5mSalesRecords.csv
  Time (mean ± σ):      1.537 s ±  0.012 s    [User: 1.548 s, System: 0.240 s]
  Range (min … max):    1.520 s …  1.566 s    10 runs
 
Benchmark 2: wc -w ./5mSalesRecords.csv
  Time (mean ± σ):      2.041 s ±  0.011 s    [User: 1.941 s, System: 0.097 s]
  Range (min … max):    2.029 s …  2.063 s    10 runs
 
Summary
  './gowc -w -bs 1000000 ./5mSalesRecords.csv' ran
    1.33 ± 0.01 times faster than 'wc -w ./5mSalesRecords.csv'

# Characters only count
$ hyperfine  --warmup 3 './gowc -m -bs 1000000 ./5mSalesRecords.csv' 'wc -m ./5mSalesRecords.csv'                                                                        [±main ●]
Benchmark 1: ./gowc -m -bs 1000000 ./5mSalesRecords.csv
  Time (mean ± σ):     751.9 ms ±   6.4 ms    [User: 707.1 ms, System: 149.5 ms]
  Range (min … max):   741.9 ms … 759.5 ms    10 runs
 
Benchmark 2: wc -m ./5mSalesRecords.csv
  Time (mean ± σ):      5.667 s ±  0.094 s    [User: 5.539 s, System: 0.113 s]
  Range (min … max):    5.578 s …  5.794 s    10 runs
 
Summary
  './gowc -m -bs 1000000 ./5mSalesRecords.csv' ran
    7.54 ± 0.14 times faster than 'wc -m ./5mSalesRecords.csv'

# Multiple files
$ hyperfine  --warmup 3 './gowc -bs 1000000 ./5mSalesRecords.csv ./5mSalesRecords.csv' 'wc ./5mSalesRecords.csv ./5mSalesRecords.csv'                [±main ●]
Benchmark 1: ./gowc -bs 1000000 ./5mSalesRecords.csv ./5mSalesRecords.csv
  Time (mean ± σ):      1.698 s ±  0.009 s    [User: 3.271 s, System: 0.515 s]
  Range (min … max):    1.684 s …  1.708 s    10 runs
 
Benchmark 2: wc ./5mSalesRecords.csv ./5mSalesRecords.csv
  Time (mean ± σ):      4.082 s ±  0.013 s    [User: 3.886 s, System: 0.192 s]
  Range (min … max):    4.062 s …  4.102 s    10 runs
 
Summary
  './gowc -bs 1000000 ./5mSalesRecords.csv ./5mSalesRecords.csv' ran
    2.40 ± 0.01 times faster than 'wc ./5mSalesRecords.csv ./5mSalesRecords.csv
```

## Tests

Run the test suite.

```bash
make test
```