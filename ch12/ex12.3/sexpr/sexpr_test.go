package sexpr

import (
	"testing"
)

func TestBool(t *testing.T) {
	type iface interface{}

	tests := []struct {
		v    interface{}
		want string
	}{
		{[2]bool{true, false}, "(t nil)"},
		{struct {
			a float32
			b float64
		}{-123, 1.234e+10}, "((a -123) (b 1.234e+10))"},
		{struct {
			a complex64
			b complex128
		}{-1 - 2i, 1.23e10 + 1.0e8i}, "((a #C(-1 -2)) (b #C(1.23e+10 1e+08)))"},
		{struct {
			v iface
		}{[]int{1, 2, 3}}, `((v ("[]int" (1 2 3))))`},
	}

	for _, test := range tests {
		b, err := Marshal(test.v)
		if err != nil {
			t.Errorf("Marshal(%v) failed: %s", test.v, err)
		}
		s := string(b)
		if s != test.want {
			t.Errorf("Marshal(%v) got %s, want %s", test.v, s, test.want)
		}
	}
}
