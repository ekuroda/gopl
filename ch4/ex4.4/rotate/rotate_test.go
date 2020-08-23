package rotate

import "testing"

func TestRotate(t *testing.T) {

	s := []int{1, 2, 3, 4, 5}
	copy := append(make([]int, 0), s...)
	Rotate(s, 2)
	want := []int{3, 4, 5, 1, 2}
	if len(s) != len(want) {
		t.Fatalf("Rotate(%v); length got %d, want %d", copy, len(s), len(want))
	}

	for i, v := range s {
		if v != want[i] {
			t.Errorf("Rotate(%v); got %v, want %v", copy, s, want)
			break
		}
	}
}
