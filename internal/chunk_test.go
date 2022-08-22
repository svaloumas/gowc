package gowc

import (
	"io/ioutil"
	"reflect"
	"testing"
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
	bufferSize := 200

	remainder := fileSize % bufferSize

	counters := []Counter{}
	c1 := Counter{
		chars:          200,
		lines:          3,
		words:          30,
		newLineIdxs:    []int{0, 57, 123},
		startsWithChar: false,
		endsWithChar:   true,
	}
	c2 := Counter{
		chars:          200,
		lines:          2,
		words:          31,
		newLineIdxs:    []int{229, 321},
		startsWithChar: true,
		endsWithChar:   false,
	}
	c3 := Counter{
		chars:          200,
		lines:          2,
		words:          31,
		newLineIdxs:    []int{410, 527},
		startsWithChar: true,
		endsWithChar:   true,
	}
	c4 := Counter{
		chars:          200,
		lines:          3,
		words:          29,
		newLineIdxs:    []int{639, 640, 713},
		startsWithChar: true,
		endsWithChar:   false,
	}
	c5 := Counter{
		chars:          199,
		lines:          3,
		words:          29,
		newLineIdxs:    []int{818, 926, 998},
		startsWithChar: true,
		endsWithChar:   false,
	}
	counters = append(counters, c1, c2, c3, c4, c5)

	opts := &Options{
		Bytes:      true,
		Words:      true,
		Lines:      true,
		Chars:      true,
		MaxLine:    true,
		BufferSize: bufferSize,
	}

	// Reset cursor.
	fp.Seek(0, 0)

	chunks := ReadFileInChunks(fp, fileSize, opts)

	i := 0
	for chunk := range chunks {

		if eq := reflect.DeepEqual(chunk.Counter, counters[i]); !eq {
			t.Errorf("ReadFileInChunks returned wrong bytes count: got %v want %v", chunk.Counter, counters[i])
		}

		bufSize := chunk.bufSize
		// Last chunk should have a buffer size equal to remainder.
		if i == len(counters)-1 {
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
		i++
	}
}
