package main

import (
	"testing"
)

func TestWordCounter(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"aaa, bbb,\nccc ddd.\neee fff ggg.", 7},
		{"", 0},
		{"aaa\t bbb\n\tccc", 3},
	}

	for _, test := range tests {
		var c WordCounter
		_, err := c.Write([]byte(test.input))
		if err != nil {
			t.Fatal(err)
		}

		if c != WordCounter(test.want) {
			t.Errorf("%q; c = %d, want %d", test.input, c, test.want)
		}
	}

	for _, test := range tests {
		var c WordCounter2
		_, err := c.Write([]byte(test.input))
		if err != nil {
			t.Fatal(err)
		}

		if c != WordCounter2(test.want) {
			t.Errorf("%q; c = %d, want %d", test.input, c, test.want)
		}
	}
}

func TestLineCounter(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"aaa, bbb,\nccc ddd.\neee fff ggg.", 3},
		{"", 0},
		{"\n", 1},
		{" ", 1},
		{" \n", 1},
		{"\n\n", 2},
		{"\n\n ", 3},
	}

	for _, test := range tests {
		var c LineCounter
		_, err := c.Write([]byte(test.input))
		if err != nil {
			t.Fatal(err)
		}

		if c != LineCounter(test.want) {
			t.Errorf("%q; c = %d, want %d", test.input, c, test.want)
		}
	}

	for _, test := range tests {
		var c LineCounter2
		_, err := c.Write([]byte(test.input))
		if err != nil {
			t.Fatal(err)
		}

		if c != LineCounter2(test.want) {
			t.Errorf("%q; c = %d, want %d", test.input, c, test.want)
		}
	}
}
