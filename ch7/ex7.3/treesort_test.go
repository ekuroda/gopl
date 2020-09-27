package main

import (
	"fmt"
	"testing"
)

func TestTreeSort(t *testing.T) {
	tests := []struct {
		numbers []int
	}{
		{[]int{5, 2, 3, 4, 1}},
		{},
	}

	for _, test := range tests {
		numbers := make([]int, len(test.numbers))
		copy(numbers, test.numbers)

		tree := sort(numbers)
		treeStr := fmt.Sprintf("%v", tree)
		numbersStr := fmt.Sprintf("%v", numbers)

		if treeStr != numbersStr {
			t.Errorf("sort(%v); tree = %q, want %q", test.numbers, treeStr, numbersStr)
		}
	}
}
