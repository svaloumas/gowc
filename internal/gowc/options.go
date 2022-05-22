package gowc

import (
	"flag"
	"fmt"
)

const Version = "v0.1.0"

const usage = `
	gowc ` + Version + `
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
			-bs --buffer-size        Configure the buffer size of each chunk to be processed.

	OPTIONS:
			--files0-from <file>    Read input from the NUL-terminated list of filenames in the given file.
			--files-from <file>     Read input from the newline-terminated list of filenames in the given file.
	ARGS:
		<input>...    Input file names
`

// Options is the flags and options provided
// by the user via the command line.
type Options struct {
	Bytes        bool
	Chars        bool
	Lines        bool
	MaxLine      bool
	Words        bool
	Version      bool
	BufferSize   int
	FilesFrom    string
	FilesFromNUL string
	Filepaths    []string
}

// ParseOptions parses the command line flags and options and stores their values in memory.
func ParseOptions() *Options {
	opts := new(Options)
	flag.BoolVar(&opts.Bytes, "c", false, "print bytes")
	flag.BoolVar(&opts.Bytes, "bytes", false, "print bytes")
	flag.BoolVar(&opts.Chars, "m", false, "print utf-8 characters")
	flag.BoolVar(&opts.Chars, "chars", false, "print utf-8 characters")
	flag.BoolVar(&opts.Lines, "l", false, "print utf-8 characters")
	flag.BoolVar(&opts.Lines, "lines", false, "print lines")
	flag.BoolVar(&opts.MaxLine, "L", false, "print the length of the longest line")
	flag.BoolVar(&opts.MaxLine, "max-line-length", false, "print the length of the longest line")
	flag.BoolVar(&opts.Words, "w", false, "print words")
	flag.BoolVar(&opts.Words, "words", false, "print words")
	flag.BoolVar(&opts.Version, "v", false, "print version")
	flag.BoolVar(&opts.Version, "version", false, "print version")
	flag.IntVar(&opts.BufferSize, "bs", 4*1024, "buffer size to process concurrently")
	flag.IntVar(&opts.BufferSize, "buffer-size", 4*1024, "buffer size to process concurrently")

	flag.StringVar(&opts.FilesFrom, "files-from", "", "read input from given file (newline terminated list of files)")
	flag.StringVar(&opts.FilesFromNUL, "files0-from", "", "read input from given file (NUL terminated list of files)")

	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	if !opts.Lines && !opts.Words && !opts.Bytes && !opts.Chars && !opts.MaxLine {
		opts.Lines = true
		opts.Words = true
		opts.Bytes = true
	}

	opts.Filepaths = flag.Args()
	return opts
}
