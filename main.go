package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"unicode"
)

type chunk struct {
	bufChan chan []byte
	bufSize int
	offset  int
}

type counter struct {
	bytes          int
	chars          int
	lines          int
	words          int
	newLineIdxs    []int
	startsWithChar bool
	endsWithChar   bool
}

func readFileInChunks(fp *os.File, fileSize, bufferSize int) []*chunk {

	concurrency := fileSize / bufferSize
	chunks := make([]*chunk, concurrency)
	for i := 0; i < concurrency; i++ {
		c := &chunk{
			bufChan: make(chan []byte),
			bufSize: bufferSize,
			offset:  bufferSize * i,
		}
		chunks[i] = c
	}

	if remainder := fileSize % bufferSize; remainder != 0 {
		c := &chunk{
			bufChan: make(chan []byte),
			bufSize: remainder,
			offset:  bufferSize * concurrency,
		}
		concurrency++
		chunks = append(chunks, c)
	}

	for i := 0; i < concurrency; i++ {
		idx := i

		go func(chunks []*chunk, i int) {

			chunk := chunks[idx]
			buf := make([]byte, chunk.bufSize)
			_, err := fp.ReadAt(buf, int64(chunk.offset))
			if err != nil {
				return
			}
			chunk.bufChan <- buf
			close(chunk.bufChan)

		}(chunks, idx)

	}
	return chunks
}

func processChunks(chunks []*chunk, opts *options) <-chan *counter {
	out := make(chan *counter)
	go func() {
		for i, chunk := range chunks {
			c := processBuffer(chunk.bufChan, opts, i)
			out <- c
		}
		close(out)
	}()
	return out
}

func processBuffer(bufChan <-chan []byte, opts *options, order int) *counter {
	bytes := <-bufChan

	bufString := string(bytes)

	chunkCounter := &counter{}
	if opts.bytes {
		chunkCounter.bytes = len(bytes)
	}
	if opts.chars {
		chunkCounter.chars = len(bufString)
	}
	if opts.lines || opts.maxLine {
		newLineIdxs := continuousNewLineIndexes(bufString, order, opts.bufferSize)
		chunkCounter.lines = len(newLineIdxs)
		chunkCounter.newLineIdxs = newLineIdxs
	}
	if opts.words {
		fields := strings.Fields(bufString)
		chunkCounter.words = len(fields)

		bufRunes := []rune(bufString)
		firstRune := bufRunes[0]
		lastRune := bufRunes[len(bufRunes)-1]

		if !unicode.IsSpace(firstRune) {
			chunkCounter.startsWithChar = true
		}
		if !unicode.IsSpace(lastRune) {
			chunkCounter.endsWithChar = true
		}
	}
	return chunkCounter
}

func maxLineLength(c *counter) int {
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
	return maxLine
}

func printResults(c *counter, maxLine int, filename string, opts *options) {
	if opts.lines {
		fmt.Printf("\t%d", c.lines)
	}
	if opts.words {
		fmt.Printf("\t%d", c.words)
	}
	if opts.bytes {
		fmt.Printf("\t%d", c.bytes)
	}
	if opts.chars {
		fmt.Printf("\t%d", c.chars)
	}
	if opts.maxLine {
		fmt.Printf("\t%d", maxLine)
	}
	fmt.Printf(" %s", filename)
}

func aggregate(c *counter, chunkCounter *counter) {
	// Not the first chunk.
	if c.words != 0 {
		if chunkCounter.startsWithChar && c.endsWithChar {
			c.words -= 1
		}
	}
	c.startsWithChar = chunkCounter.startsWithChar
	c.endsWithChar = chunkCounter.endsWithChar

	c.bytes += chunkCounter.bytes
	c.chars += chunkCounter.chars
	c.lines += chunkCounter.lines
	c.words += chunkCounter.words

	c.newLineIdxs = append(c.newLineIdxs, chunkCounter.newLineIdxs...)
}

func continuousNewLineIndexes(s string, chunkOrder, bufferSize int) []int {
	indexes := []int{}
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			orderedIdx := i + (chunkOrder * bufferSize)
			indexes = append(indexes, orderedIdx)
		}
	}
	return indexes
}

func main() {
	var wg sync.WaitGroup
	var maxLine int
	opts := parseOptions()

	for _, filepath := range opts.filepaths {
		fp, err := os.Open(filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer fp.Close()
		fi, err := fp.Stat()
		if err != nil {
			log.Fatal(err)
		}

		file := filepath
		wg.Add(1)
		go func(string) {
			c := &counter{}
			chunks := readFileInChunks(fp, int(fi.Size()), opts.bufferSize)
			chunkCounterChan := processChunks(chunks, opts)

			for chunkCounter := range chunkCounterChan {
				aggregate(c, chunkCounter)
			}
			if opts.maxLine {
				maxLine = maxLineLength(c)
			}
			printResults(c, maxLine, file, opts)
			wg.Done()
		}(file)
	}
	wg.Wait()
}
