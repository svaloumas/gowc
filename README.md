# gowc
[![CI](https://github.com/svaloumas/gowc/actions/workflows/ci.yml/badge.svg)](https://github.com/svaloumas/gowc/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/svaloumas/gowc/branch/main/graph/badge.svg?token=9CI4Q74JJK)](https://codecov.io/gh/svaloumas/gowc)
[![Go Report Card](https://goreportcard.com/badge/github.com/svaloumas/gowc)](https://goreportcard.com/report/github.com/svaloumas/gowc)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/svaloumas/gowc/blob/main/LICENSE)

Just another [`wc`](https://www.gnu.org/software/coreutils/manual/html_node/wc-invocation.html#wc-invocation) clone written in Go.

## Overview

`gowc` is a simple, zero-dependency command line tool for counting bytes, characters, words and newlines in each given file.
It leverages the language's built-in support for concurrency by processing the given input files in chunks. The buffer size of each chunk is configurable
and can be set via `-bs --buffer-size` flag. The number of go-routines that process the chunks concurrently is calculated as follows `concurrency = filesize / buffersize`.

## Usage

By default, `gowc` will count lines, words, and bytes. You can specify the counters you'd like by using the available flags and options from the table below.

```
gowc v0.1.0
A Go word count (wc) clone - print newline, word, char, and byte counts for each file 

USAGE:
	gowc [FLAGS] [OPTIONS] [input]...

FLAGS:
		-c, --bytes              Print the byte counts.
		-m, --chars              Print the character counts.
		-l, --lines              Print the newline counts.
		-L, --max-line-length    Print the length of the longest line.
		-w, --words              Print the word counts.
		-h, --help               Display help and exit.
		-V, --version            Output version information and exit.
		-bs --buffer-size        Configure the buffer size of each chunk to be processed (defaults to 4096).

OPTIONS:
		--files-from <file>     Read input from the files specified by a newline-terminated list of filenames in the given file.
ARGS:
	<input>...    Input filenames
```

## Performance

TBD

## Tests

Run the test suite.

```bash
make test
```