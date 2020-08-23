package reverse

import "testing"

func TestReverse(t *testing.T) {
	test := [...]int{1, 2, 3, 4, 5}
	copy := test
	Reverse(&test)
	if test != [...]int{5, 4, 3, 2, 1} {
		t.Errorf("reverse(%v); got %v, want %v", copy, test, [...]int{5, 4, 3, 2, 1})
	}
}
