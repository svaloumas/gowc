package gowc

import (
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"unicode"
)

var testFileContent = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit.
Nunc ut mattis lectus. Maecenas congue magna sed commodo pretium.
Ut ornare, felis sed tempus venenatis, risus nunc tempor urna, tempus facilisis arcu velit sit amet arcu.
Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos.
Vestibulum faucibus orci et tincidunt tempor. Nunc non augue et ligula tempus porttitor.
Praesent pulvinar felis hendrerit lacinia mattis. Quisque ac eleifend enim. Sed et bibendum dolor, in eleifend urna.
Maecenas eu lectus id leo bibendum mollis id a ligula. Donec porta hendrerit ex, sed euismod ipsum bibendum at.

Vestibulum orci ante, pretium vitae erat vitae, gravida elementum lacus.
Donec venenatis interdum turpis ornare congue. Aenean tempor sem tortor, quis rhoncus urna convallis ut.
Praesent sagittis at lacus non dignissim. Donec ipsum nunc, accumsan eu aliquet dictum, interdum et libero.
Fusce laoreet laoreet ex quis ornare. Pellentesque vitae pharetra nibh.
`

func TestReadFileInChunks(t *testing.T) {
	dir := t.TempDir()
	fp, _ := ioutil.TempFile(dir, "test_file.txt")
	defer fp.Close()

	fp.WriteString(testFileContent)

	fi, _ := fp.Stat()
	fileSize := int(fi.Size())
	bufferSize := 30

	numOfChunks := fileSize / bufferSize
	remainder := fileSize % bufferSize

	expectedBytesChunks := [][]byte{}

	for i := 0; i < numOfChunks; i++ {
		bytesChunk := make([]byte, bufferSize)
		fp.ReadAt(bytesChunk, int64(bufferSize*i))
		expectedBytesChunks = append(expectedBytesChunks, bytesChunk)
	}
	bytesChunkRemainder := make([]byte, remainder)
	fp.ReadAt(bytesChunkRemainder, int64(bufferSize*numOfChunks))
	expectedBytesChunks = append(expectedBytesChunks, bytesChunkRemainder)

	// Reset cursor.
	fp.Seek(0, 0)

	chunks := ReadFileInChunks(fp, fileSize, bufferSize)

	for i, chunk := range chunks {
		bytes := <-chunk.bufChan

		if eq := reflect.DeepEqual(bytes, expectedBytesChunks[i]); !eq {
			t.Errorf("ReadFileInChunks returned wrong bytes chunk: got %v want %v", bytes, expectedBytesChunks[i])
		}

		bufSize := chunk.bufSize
		// Last chunk should have a buffer size equal to remainder.
		if i == len(chunks)-1 {
			if bufSize != remainder {
				t.Errorf("ReadFileInChunks returned wrong chunk buffer size: got %v want %v", bufSize, remainder)
			}
		} else {
			if bufSize != bufferSize {
				t.Errorf("ReadFileInChunks returned wrong chunk buffer size: got %v want %v", bufSize, bufferSize)
			}
		}

		offset := chunk.offset
		if offset != bufferSize*i {
			t.Errorf("ReadFileInChunks returned wrong chunk buffer size: got %v want %v", offset, bufferSize*i)
		}
	}
}

func TestProcessChunks(t *testing.T) {
	dir := t.TempDir()
	fp, _ := ioutil.TempFile(dir, "test_file.txt")
	defer fp.Close()

	fp.WriteString(testFileContent)

	fi, _ := fp.Stat()
	fileSize := int(fi.Size())
	bufferSize := 200

	numOfChunks := fileSize / bufferSize
	remainder := fileSize % bufferSize

	bytesChunks := [][]byte{}

	for i := 0; i < numOfChunks; i++ {
		bytesChunk := make([]byte, bufferSize)
		fp.ReadAt(bytesChunk, int64(bufferSize*i))
		bytesChunks = append(bytesChunks, bytesChunk)
	}
	bytesChunkRemainder := make([]byte, remainder)
	fp.ReadAt(bytesChunkRemainder, int64(bufferSize*numOfChunks))
	bytesChunks = append(bytesChunks, bytesChunkRemainder)

	opts := &Options{
		Bytes:      true,
		Chars:      true,
		Lines:      true,
		MaxLine:    true,
		Words:      true,
		BufferSize: bufferSize,
	}

	chunks := []*chunk{}
	expectedCounters := []*Counter{}
	for order, bytesChunk := range bytesChunks {
		bytes := bytesChunk
		bufSize := bufferSize
		// Last chunk.
		if order == len(bytesChunks)-1 {
			bufSize = remainder
		}
		ch := &chunk{
			bufChan: make(chan []byte),
			bufSize: bufSize,
			offset:  bufferSize * order,
		}
		chunks = append(chunks, ch)

		go func([]byte) {
			ch.bufChan <- bytes
			close(ch.bufChan)
		}(bytes)

		bufString := string(bytes)

		c := &Counter{}
		c.bytes = len(bytes)
		c.chars = len(bufString)
		newLineIdxs := []int{}
		for i := 0; i < len(bufString); i++ {
			if bufString[i] == '\n' {
				orderedIdx := i + (order * bufferSize)
				newLineIdxs = append(newLineIdxs, orderedIdx)
			}
		}
		c.lines = len(newLineIdxs)
		c.newLineIdxs = newLineIdxs
		fields := strings.Fields(bufString)
		c.words = len(fields)

		bufRunes := []rune(bufString)
		if len(bufRunes) >= 1 {

			firstRune := bufRunes[0]
			lastRune := bufRunes[len(bufRunes)-1]

			if !unicode.IsSpace(firstRune) {
				c.startsWithChar = true
			}
			if !unicode.IsSpace(lastRune) {
				c.endsWithChar = true
			}
		}
		expectedCounters = append(expectedCounters, c)
	}

	counterChan := ProcessChunks(chunks, opts)

	i := 0
	for counter := range counterChan {
		if eq := reflect.DeepEqual(counter, expectedCounters[i]); !eq {
			t.Errorf("ProcessChunks returned wrong counter: got %v want %v", counter, expectedCounters[i])
		}
		i++
	}
}
