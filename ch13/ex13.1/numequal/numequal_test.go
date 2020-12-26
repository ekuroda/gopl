package numequal

import (
	"testing"
)

func TestEqual(t *testing.T) {
	for _, test := range []struct {
		x, y interface{}
		want bool
	}{
		{-1, -1, true},
		{-1, -2, false},
		{uint(1), uint(1), true},
		{uint(1), uint(2), false},
		{0.1, 0.1, true},
		{0.0, 0.0, true},
		{0.0, -1e-10, true},
		{0.0, -1e-9, false},
		{0.0, 1e-8, false},
		{complex(0, 0), complex(0, 0), true},
		{complex(1e-10, 1e-10), complex(0, 0), true},
		{complex(1e-9, 1e-10), complex(0, 0), false},
		{complex(1e-10, 1e-9), complex(0, 0), false},
	} {
		if Equal(test.x, test.y) != test.want {
			t.Errorf("Equal(%.10f, %.10f) = %t", test.x, test.y, !test.want)
		}
	}
}
