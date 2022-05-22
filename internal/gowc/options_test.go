package gowc

import (
	"testing"
)

func TestParseOptions(t *testing.T) {
	opts := ParseOptions()

	if opts.Lines != true {
		t.Errorf("ParseOptions returned wrong lines flag value: got %v want true", opts.Lines)
	}
	if opts.Words != true {
		t.Errorf("ParseOptions returned wrong words flag value: got %v want true", opts.Words)
	}
	if opts.Bytes != true {
		t.Errorf("ParseOptions returned wrong bytes flag value: got %v want true", opts.Bytes)
	}

	if opts.Chars != false {
		t.Errorf("ParseOptions returned wrong chars flag value: got %v want false", opts.Chars)
	}
	if opts.MaxLine != false {
		t.Errorf("ParseOptions returned wrong max line flag value: got %v want false", opts.MaxLine)
	}
	if opts.Version != false {
		t.Errorf("ParseOptions returned wrong version flag value: got %v want false", opts.Version)
	}

	if opts.BufferSize != 4096 {
		t.Errorf("ParseOptions returned wrong files from flag value: got %v want 4096", opts.BufferSize)
	}
}
