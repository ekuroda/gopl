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

	// for _, test := range tests {
	// 	if got := comma(test.input); got != test.want {
	// 		t.Errorf("comma(%s) = %s, want %s", test.input, got, test.want)
	// 	}
	// }

	for _, test := range tests {
		if got := commaByLoop(test.input); got != test.want {
			t.Errorf("commaByLoop(%s) = %s; want %s", test.input, got, test.want)
		}
	}
}
