package gowc

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestMaxLineLength(t *testing.T) {
	c := &Counter{
		newLineIdxs: []int{12, 14, 35, 59, 66, 90, 101, 104},
	}
	// 90 - 66 - 1
	expected := 23

	maxLineLength := c.MaxLineLength()

	if maxLineLength != expected {
		t.Errorf("MaxLineLength returned wrong length: got %v want %v", maxLineLength, expected)
	}
}

func TestAggregate(t *testing.T) {
	firstChunkCounter := Counter{
		chars:          1000,
		lines:          3,
		words:          100,
		newLineIdxs:    []int{3, 30, 90},
		startsWithChar: true,
		endsWithChar:   true,
	}
	secondChunkCounter := Counter{
		chars:          1000,
		lines:          3,
		words:          100,
		newLineIdxs:    []int{103, 130, 190},
		startsWithChar: true,
		endsWithChar:   true,
	}
	expected := &Counter{
		chars:          2000,
		lines:          6,
		words:          199,
		newLineIdxs:    []int{3, 30, 90, 103, 130, 190},
		startsWithChar: true,
		endsWithChar:   true,
	}

	c := &Counter{}
	c.Aggregate(firstChunkCounter)
	c.Aggregate(secondChunkCounter)

	if eq := reflect.DeepEqual(c, expected); !eq {
		t.Errorf("Aggregate returned wrong aggregate counter: got %v want %v", c, expected)
	}
}

func TestPrintCounter(t *testing.T) {
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	c := &Counter{
		Bytes:          2000,
		chars:          2000,
		lines:          6,
		words:          199,
		newLineIdxs:    []int{3, 30, 90, 103, 130, 190},
		startsWithChar: true,
		endsWithChar:   true,
	}

	opts := &Options{
		Bytes:   true,
		Chars:   true,
		Lines:   true,
		MaxLine: true,
		Words:   true,
	}

	filename := "test_file.txt"
	c.PrintCounter(100, filename, opts)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = stdout

	expected := fmt.Sprintf("\t6\t199\t2000\t2000\t100 %s\n", filename)

	if string(out) != expected {
		t.Errorf("PrintCounter printed wrong output: got %s want %s", string(out), expected)
	}
}
