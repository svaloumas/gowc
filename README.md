# gowc
[![CI](https://github.com/svaloumas/gowc/actions/workflows/ci.yml/badge.svg)](https://github.com/svaloumas/gowc/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/svaloumas/gowc/branch/main/graph/badge.svg?token=7DvuWdQZPr)](https://codecov.io/gh/svaloumas/gowc)
[![Go Report Card](https://goreportcard.com/badge/github.com/svaloumas/gowc)](https://goreportcard.com/report/github.com/svaloumas/gowc)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/svaloumas/gowc/blob/main/LICENSE)

Just another GNU [`wc`](https://www.gnu.org/software/coreutils/manual/html_node/wc-invocation.html#wc-invocation) clone, written in Go.

## Overview

`gowc` is a simple, zero-dependency command line tool for counting bytes, characters, words and newlines in each given file.
It leverages the language's built-in support for concurrency by processing the given input files in chunks. The buffer size of each chunk is configurable
and can be set via `-bs, --buffer-size` flag. The number of go-routines that process the chunks concurrently is calculated as follows `concurrency = filesize / buffersize`.

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
$ hyperfine  --warmup 3 './gowc -l -bs 1000000 ./5mSalesRecords.csv' 'wc -l ./5mSalesRecords.csv'
Benchmark 1: ./gowc -l -bs 1000000 ./5mSalesRecords.csv
  Time (mean ± σ):     145.9 ms ±  14.3 ms    [User: 290.0 ms, System: 612.6 ms]
  Range (min … max):   121.5 ms … 170.0 ms    20 runs
 
Benchmark 2: wc -l ./5mSalesRecords.csv
  Time (mean ± σ):     472.3 ms ±   6.1 ms    [User: 384.5 ms, System: 86.3 ms]
  Range (min … max):   467.7 ms … 488.4 ms    10 runs
 
Summary
  './gowc -l -bs 1000000 ./5mSalesRecords.csv' ran
    3.24 ± 0.32 times faster than 'wc -l ./5mSalesRecords.csv'

# Default lines, words and bytes count
$ hyperfine  --warmup 3 './gowc -bs 1000000 ./5mSalesRecords.csv' 'wc ./5mSalesRecords.csv'
Benchmark 1: ./gowc -bs 1000000 ./5mSalesRecords.csv
  Time (mean ± σ):     444.2 ms ±  17.5 ms    [User: 2571.0 ms, System: 493.4 ms]
  Range (min … max):   423.8 ms … 480.8 ms    10 runs
 
Benchmark 2: wc ./5mSalesRecords.csv
  Time (mean ± σ):      2.020 s ±  0.009 s    [User: 1.925 s, System: 0.092 s]
  Range (min … max):    2.009 s …  2.035 s    10 runs
 
Summary
  './gowc -bs 1000000 ./5mSalesRecords.csv' ran
    4.55 ± 0.18 times faster than 'wc ./5mSalesRecords.csv'

# Word only count
$ hyperfine  --warmup 3 './gowc -w -bs 1000000 ./5mSalesRecords.csv' 'wc -w ./5mSalesRecords.csv'
Benchmark 1: ./gowc -w -bs 1000000 ./5mSalesRecords.csv
  Time (mean ± σ):     429.8 ms ±  18.7 ms    [User: 2500.1 ms, System: 474.4 ms]
  Range (min … max):   409.7 ms … 464.2 ms    10 runs
 
Benchmark 2: wc -w ./5mSalesRecords.csv
  Time (mean ± σ):      2.004 s ±  0.010 s    [User: 1.912 s, System: 0.090 s]
  Range (min … max):    1.991 s …  2.022 s    10 runs
 
Summary
  './gowc -w -bs 1000000 ./5mSalesRecords.csv' ran
    4.66 ± 0.20 times faster than 'wc -w ./5mSalesRecords.csv'

# Characters only count
$ hyperfine  --warmup 3 './gowc -m -bs 1000000 ./5mSalesRecords.csv' 'wc -m ./5mSalesRecords.csv'
Benchmark 1: ./gowc -m -bs 1000000 ./5mSalesRecords.csv
  Time (mean ± σ):     241.2 ms ±   9.1 ms    [User: 1157.2 ms, System: 450.7 ms]
  Range (min … max):   229.8 ms … 263.6 ms    12 runs
 
Benchmark 2: wc -m ./5mSalesRecords.csv
  Time (mean ± σ):      5.467 s ±  0.014 s    [User: 5.364 s, System: 0.097 s]
  Range (min … max):    5.451 s …  5.501 s    10 runs
 
Summary
  './gowc -m -bs 1000000 ./5mSalesRecords.csv' ran
   22.66 ± 0.85 times faster than 'wc -m ./5mSalesRecords.csv'

# Multiple files
$ hyperfine  --warmup 3 './gowc -bs 1000000 ./5mSalesRecords.csv ./5mSalesRecords.csv' 'wc ./5mSalesRecords.csv ./5mSalesRecords.csv'
Benchmark 1: ./gowc -bs 1000000 ./5mSalesRecords.csv ./5mSalesRecords.csv
  Time (mean ± σ):     849.4 ms ±  33.1 ms    [User: 5082.0 ms, System: 844.0 ms]
  Range (min … max):   816.1 ms … 929.3 ms    10 runs
 
Benchmark 2: wc ./5mSalesRecords.csv ./5mSalesRecords.csv
  Time (mean ± σ):      4.205 s ±  0.197 s    [User: 3.975 s, System: 0.205 s]
  Range (min … max):    3.951 s …  4.502 s    10 runs
 
Summary
  './gowc -bs 1000000 ./5mSalesRecords.csv ./5mSalesRecords.csv' ran
    4.95 ± 0.30 times faster than 'wc ./5mSalesRecords.csv ./5mSalesRecords.csv'
```

## Tests

Run the test suite.

```bash
make test
```