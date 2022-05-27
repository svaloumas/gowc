package gowc

import (
	"fmt"
)

// Counter is wrapper for all the metrics.
type Counter struct {
	Bytes          int
	chars          int
	lines          int
	words          int
	newLineIdxs    []int
	startsWithChar bool
	endsWithChar   bool
}

// MaxLineLength calculates the longest line length for an aggregated counter.
func (c *Counter) MaxLineLength() int {
	var maxLine int
	if len(c.newLineIdxs) > 0 {
		for i, idx := range c.newLineIdxs {
			if i == 0 {
				maxLine = c.newLineIdxs[i]
				continue
			}
			// Subtrack the current position of idx.
			nextLine := idx - c.newLineIdxs[i-1] - 1
			if nextLine > maxLine {
				maxLine = nextLine
			}
		}
	}
	lastLine := c.Bytes - c.newLineIdxs[len(c.newLineIdxs)-1] - 1
	if lastLine > maxLine {
		maxLine = lastLine
	}
	return maxLine
}

// Aggregate aggregates all the metrics from the chunk counters into a main file counter.
func (c *Counter) Aggregate(chunkCounter Counter) {
	// Not the first chunk.
	if c.words != 0 {
		if chunkCounter.startsWithChar && c.endsWithChar {
			c.words -= 1
		}
	}
	c.startsWithChar = chunkCounter.startsWithChar
	c.endsWithChar = chunkCounter.endsWithChar

	c.chars += chunkCounter.chars
	c.lines += chunkCounter.lines
	c.words += chunkCounter.words

	c.newLineIdxs = append(c.newLineIdxs, chunkCounter.newLineIdxs...)
}

// PrintCounter outputs the counter results.
func (c *Counter) PrintCounter(maxLine int, filename string, opts *Options) {
	if opts.Lines {
		fmt.Printf("\t%d", c.lines)
	}
	if opts.Words {
		fmt.Printf("\t%d", c.words)
	}
	if opts.Bytes {
		fmt.Printf("\t%d", c.Bytes)
	}
	if opts.Chars {
		fmt.Printf("\t%d", c.chars)
	}
	if opts.MaxLine {
		fmt.Printf("\t%d", maxLine)
	}
	fmt.Printf(" %s\n", filename)
}
