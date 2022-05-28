package gowc

import (
	"bytes"
	"os"
	"unicode"
	"unicode/utf8"
)

type Chunk struct {
	CounterChan chan Counter
	bufSize     int
	offset      int
}

// ReadFileInChunks reads a given file in chunks and returns a slice of chunks.
func ReadFileInChunks(fp *os.File, fileSize int, opts *Options) []*Chunk {

	concurrency := fileSize / opts.BufferSize
	chunks := make([]*Chunk, concurrency)
	for i := 0; i < concurrency; i++ {
		c := &Chunk{
			CounterChan: make(chan Counter),
			bufSize:     opts.BufferSize,
			offset:      opts.BufferSize * i,
		}
		chunks[i] = c
	}

	remainder := fileSize % opts.BufferSize
	if remainder != 0 {
		c := &Chunk{
			CounterChan: make(chan Counter),
			bufSize:     remainder,
			offset:      opts.BufferSize * concurrency,
		}
		concurrency++
		chunks = append(chunks, c)
	}

	for i := 0; i < concurrency; i++ {
		idx := i

		go func(chunks []*Chunk, idx int) {

			chunk := chunks[idx]
			buf := make([]byte, chunk.bufSize)
			_, err := fp.ReadAt(buf, int64(chunk.offset))
			if err != nil {
				return
			}

			go func() {
				chunkCounter := processBuffer(buf, opts, idx)
				chunk.CounterChan <- chunkCounter
				close(chunk.CounterChan)
			}()

		}(chunks, idx)
	}
	return chunks
}

func processBuffer(byteArr []byte, opts *Options, order int) Counter {
	chunkCounter := Counter{}

	if opts.Chars {
		chunkCounter.chars = utf8.RuneCount(byteArr)
	}

	if opts.MaxLine {
		newLineIdxs := continuousNewLineIndexes(byteArr, order, opts.BufferSize)
		chunkCounter.newLineIdxs = newLineIdxs
	}

	if opts.Lines {
		chunkCounter.lines = bytes.Count(byteArr, []byte{'\n'})
	}

	if opts.Words {
		chunkCounter.words = len(bytes.Fields(byteArr))

		firstRune := []rune(string(byteArr[0:1]))
		lastRune := []rune(string(byteArr[len(byteArr)-1:]))
		if len(firstRune) >= 1 && !unicode.IsSpace(firstRune[0]) {
			chunkCounter.startsWithChar = true
		}
		if len(lastRune) >= 1 && !unicode.IsSpace(lastRune[0]) {
			chunkCounter.endsWithChar = true
		}
	}
	return chunkCounter
}

func continuousNewLineIndexes(byteArr []byte, chunkOrder, bufferSize int) []int {
	indexes := []int{}
	placement := (chunkOrder * bufferSize)
	for i, b := range byteArr {
		if b == '\n' {
			orderedIdx := i + placement
			indexes = append(indexes, orderedIdx)
		}
	}
	return indexes
}
