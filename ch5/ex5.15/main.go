package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	var vals []int
	args := os.Args
	for i := 1; i < len(args); i++ {
		val, err := strconv.Atoi(args[i])
		if err != nil {
			log.Fatalf("invalid argument %s: %s", args[i], err)
		}
		vals = append(vals, val)
	}

	max, err := max(vals...)
	if err != nil {
		log.Fatalf("failed to calc max: %s", err)
	}
	fmt.Printf("max: %d\n", max)

	min, err := min(vals...)
	if err != nil {
		log.Fatalf("failed to calc min: %s", err)
	}
	fmt.Printf("min: %d\n", min)
}

func sum(vals ...int) int {
	total := 0
	for _, val := range vals {
		total += val
	}
	return total
}

func max(vals ...int) (int, error) {
	if len(vals) == 0 {
		return 0, fmt.Errorf("invalid arguments")
	}

	max := math.MinInt32
	for _, val := range vals {
		if val > max {
			max = val
		}
	}
	return max, nil
}

func min(vals ...int) (int, error) {
	if len(vals) == 0 {
		return 0, fmt.Errorf("invalid arguments")
	}

	min := math.MaxInt32
	for _, val := range vals {
		if val < min {
			min = val
		}
	}
	return min, nil
}

func max2(first int, vals ...int) int {
	max := first
	for _, val := range vals {
		if val > max {
			max = val
		}
	}
	return max
}

func min2(first int, vals ...int) int {
	min := first
	for _, val := range vals {
		if val < min {
			min = val
		}
	}
	return min
}
