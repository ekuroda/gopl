package main

import "testing"

func TestComma(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"1", "1"},
		{"123456", "123,456"},
		{"12345678", "12,345,678"},
	}

	for _, test := range tests {
		if got := comma(test.input); got != test.want {
			t.Errorf("comma(%s) = %s, want %s", test.input, got, test.want)
		}
	}

	for _, test := range tests {
		if got := commaByLoop(test.input); got != test.want {
			t.Errorf("commaByLoop(%s) = %s; want %s", test.input, got, test.want)
		}
	}
}

func TestCommaFloatingPointNumberWithSign(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"1", "1"},
		{"123456", "123,456"},
		{"12345678", "12,345,678"},
		{"1.1234", "1.1234"},
		{"123456.1", "123,456.1"},
		{"12345.123", "12,345.123"},
		{"12345678.12345", "12,345,678.12345"},
		{"+1234567.12", "+1,234,567.12"},
		{"-1234.12", "-1,234.12"},
	}

	for _, test := range tests {
		if got := commaFloatingPointNumberWithSign(test.input); got != test.want {
			t.Errorf("commaFloatingPointNumberWithSign(%s) = %s; want %s", test.input, got, test.want)
		}
	}
}
