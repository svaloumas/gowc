package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"unicode"
)

// Should be 4 * 1024 when shipped
const BufferSize = 10

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

func readFileInChunks(fp *os.File, fileSize int) []chunk {

	concurrency := fileSize / BufferSize
	chunks := make([]chunk, concurrency)
	for i := 0; i < concurrency; i++ {
		chunks[i].bufChan = make(chan []byte)
		chunks[i].bufSize = BufferSize
		chunks[i].offset = BufferSize * i
	}

	if remainder := fileSize % BufferSize; remainder != 0 {
		c := chunk{
			bufChan: make(chan []byte),
			bufSize: remainder,
			offset:  BufferSize * concurrency,
		}
		concurrency++
		chunks = append(chunks, c)
	}

	for i := 0; i < concurrency; i++ {
		idx := i

		go func(chunks []chunk, i int) {

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

func processChunks(chunks []chunk) <-chan *counter {
	out := make(chan *counter)
	go func() {
		for i, chunk := range chunks {
			c := processBuffer(chunk.bufChan, i)
			out <- c
		}
		close(out)
	}()
	return out
}

func processBuffer(bufChan <-chan []byte, order int) *counter {
	bytes := <-bufChan

	bufString := string(bytes)
	fields := strings.Fields(bufString)

	bufRunes := []rune(bufString)
	firstRune := bufRunes[0]
	lastRune := bufRunes[len(bufRunes)-1]

	newLineIdxs := continuousNewLineIndexes(bufString, order)
	chunkCounter := &counter{
		bytes:       len(bytes),
		chars:       len(bufString),
		lines:       len(newLineIdxs),
		words:       len(fields),
		newLineIdxs: newLineIdxs,
	}
	if !unicode.IsSpace(firstRune) {
		chunkCounter.startsWithChar = true
	}
	if !unicode.IsSpace(lastRune) {
		chunkCounter.endsWithChar = true
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

func aggregate(c *counter, chunkCounter *counter) {
	// Not the first chunk.
	if c.bytes != 0 {
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

func continuousNewLineIndexes(s string, chunkOrder int) []int {
	indexes := []int{}
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			orderedIdx := i + (chunkOrder * BufferSize)
			indexes = append(indexes, orderedIdx)
		}
	}
	return indexes
}

func main() {
	var wg sync.WaitGroup
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
			bytesChan := readFileInChunks(fp, int(fi.Size()))
			chunkCounterChan := processChunks(bytesChan)

			for chunkCounter := range chunkCounterChan {
				aggregate(c, chunkCounter)
			}
			maxLine := maxLineLength(c)

			fmt.Println(c.bytes)
			fmt.Println(c.chars)
			fmt.Println(c.lines)
			fmt.Println(c.words)
			fmt.Println(maxLine)
			fmt.Println(file)
			wg.Done()
		}(file)
	}
	wg.Wait()
}
