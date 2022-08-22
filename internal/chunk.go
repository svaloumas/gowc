package gowc

import (
	"bytes"
	"log"
	"os"
	"unicode"
	"unicode/utf8"
)

type Chunk struct {
	Counter Counter
	bufSize int
	offset  int
}

// ReadFileInChunks reads a given file in chunks and returns a slice of chunks.
func ReadFileInChunks(fp *os.File, fileSize int, opts *Options) chan *Chunk {

	numOfChunks := fileSize / opts.BufferSize
	chunks := make(chan *Chunk, 1)
	go func() {
		for i := 0; i < numOfChunks; i++ {
			c := &Chunk{
				Counter: Counter{},
				bufSize: opts.BufferSize,
				offset:  opts.BufferSize * i,
			}

			buf := make([]byte, c.bufSize)
			_, err := fp.ReadAt(buf, int64(c.offset))
			if err != nil {
				log.Printf("error reading file: %s", err)
				return
			}

			c.Counter = processBuffer(buf, opts, i)
			chunks <- c
		}

		remainder := fileSize % opts.BufferSize
		if remainder != 0 {
			c := &Chunk{
				Counter: Counter{},
				bufSize: remainder,
				offset:  opts.BufferSize * numOfChunks,
			}

			buf := make([]byte, c.bufSize)
			_, err := fp.ReadAt(buf, int64(c.offset))
			if err != nil {
				log.Printf("error reading file: %s", err)
				return
			}

			c.Counter = processBuffer(buf, opts, numOfChunks)
			chunks <- c
		}
		close(chunks)
	}()
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
