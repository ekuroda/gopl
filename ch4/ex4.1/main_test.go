package main

import "testing"

func TestBitDiffCount(t *testing.T) {
	tests := []struct {
		a    [32]uint8
		b    [32]uint8
		want int
	}{
		{
			a:    [32]uint8{0, 1, 2, 3}, // 000 001 010 011
			b:    [32]uint8{1, 2, 3, 4}, // 001 010 011 100
			want: 7,
		},
	}

	for _, test := range tests {
		got := bitDiffCount(test.a, test.b)
		if got != test.want {
			t.Errorf("bitDiffCount(%v, %v) = %d; want %d", test.a, test.b, got, test.want)
		}
	}
}
