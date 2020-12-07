package sexpr

import (
	"testing"
)

func Test(t *testing.T) {
	type A struct {
		Bool            bool
		privateBool     bool
		Int             int
		Float           float64
		String          string
		Array           [2]int
		privateArray    [1]int
		Slice           []string
		SlicePtr        *[]string
		privateSlicePtr *[]string
		Ptr             *string
		privatePtr      *int
	}
	a := A{}

	data, err := Marshal(a)
	if err != nil {
		t.Fatalf("Marshal(%v) failed: %v", a, err)
	}
	want := "()"
	if string(data) != want {
		t.Errorf("Marshal(%v) = %q, want %q\n", a, data, want)
	}
}
