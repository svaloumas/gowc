package main

import (
	"flag"
	"fmt"
)

const usage = `
gowc 0.1.0
A Go word count (wc) clone - print newline, word, and byte counts for each file 

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

OPTIONS:
		--files0-from <file>    Read input from the NUL-terminated list of filenames in the given file.
		--files-from <file>     Read input from the newline-terminated list of filenames in the given file.
ARGS:
	<input>...    Input file names
`

type options struct {
	bytes     bool
	chars     bool
	lines     bool
	maxLine   bool
	words     bool
	version   bool
	filesFrom string
	filepaths []string
}

func parseOptions() *options {
	opts := new(options)
	flag.BoolVar(&opts.bytes, "c", false, "print bytes")
	flag.BoolVar(&opts.bytes, "bytes", false, "print bytes")
	flag.BoolVar(&opts.chars, "m", false, "print utf-8 characters")
	flag.BoolVar(&opts.chars, "chars", false, "print utf-8 characters")
	flag.BoolVar(&opts.lines, "l", false, "print utf-8 characters")
	flag.BoolVar(&opts.lines, "lines", false, "print lines")
	flag.BoolVar(&opts.maxLine, "L", false, "print the length of the longest line")
	flag.BoolVar(&opts.maxLine, "max-line-length", false, "print the length of the longest line")
	flag.BoolVar(&opts.words, "w", false, "print words")
	flag.BoolVar(&opts.words, "words", false, "print words")
	flag.BoolVar(&opts.version, "v", false, "print version")
	flag.BoolVar(&opts.version, "version", false, "print version")

	flag.StringVar(&opts.filesFrom, "files0-from", "", "read input from given file")

	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	opts.filepaths = flag.Args()
	return opts
}
