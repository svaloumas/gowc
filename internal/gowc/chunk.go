package gowc

import (
	"os"
	"strings"
	"unicode"
)

type chunk struct {
	bufChan chan []byte
	bufSize int
	offset  int
}

// ReadFileInChunks reads a given file in chunks and returns a slice of chunks.
func ReadFileInChunks(fp *os.File, fileSize, bufferSize int) []*chunk {

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

// ProcessChunks process the in-memory chucks in order and produces a counter channel,
// for the consumer to wait the results for each chunk on.
func ProcessChunks(chunks []*chunk, opts *Options) <-chan *Counter {
	out := make(chan *Counter)
	go func() {
		for i, chunk := range chunks {
			c := processBuffer(chunk.bufChan, opts, i)
			out <- c
		}
		close(out)
	}()
	return out
}

func processBuffer(bufChan <-chan []byte, opts *Options, order int) *Counter {
	bytes := <-bufChan

	bufString := string(bytes)

	chunkCounter := &Counter{}
	if opts.Bytes {
		chunkCounter.bytes = len(bytes)
	}
	if opts.Chars {
		chunkCounter.chars = len(bufString)
	}
	if opts.Lines || opts.MaxLine {
		newLineIdxs := continuousNewLineIndexes(bufString, order, opts.BufferSize)
		chunkCounter.lines = len(newLineIdxs)
		chunkCounter.newLineIdxs = newLineIdxs
	}
	if opts.Words {
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
