package main

import (
	"bytes"
	"testing"
)

func TestCountingWriter(t *testing.T) {
	texts := []string{
		"0123456789",
		"abcd",
		"",
		"efg",
	}

	buf := bytes.NewBuffer(nil)
	cw, count := CountingWriter(buf)

	var total int64
	for _, text := range texts {
		bytes := []byte(text)
		_, err := cw.Write(bytes)
		if err != nil {
			t.Fatalf("failed to write: %v", err)
		}
		total += int64(len(bytes))

		if total != *count {
			t.Errorf("cw.Writer(%q); count = %d, want %d", text, *count, total)
		}
	}
}
