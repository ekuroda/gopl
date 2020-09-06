package main

import "testing"

func TestMax(t *testing.T) {
	vals := []int{3, 5, -1, -4, 2}
	m, err := max(vals...)
	if err != nil {
		t.Errorf("err of max(%v) = nil; want error", vals)
	}
	if m != 5 {
		t.Errorf("max(%v) = %d; want %d", vals, m, 5)
	}

	vals = []int{}
	_, err = max(vals...)
	if err == nil {
		t.Errorf("err of max(%v) = nil; want error", vals)
	}
}

func TestMin(t *testing.T) {
	vals := []int{3, 5, -1, -4, 2}
	m, err := min(vals...)
	if err != nil {
		t.Errorf("err of min(%v) = nil; want error", vals)
	}
	if m != -4 {
		t.Errorf("min(%v) = %d; want %d", vals, m, -4)
	}

	vals = []int{}
	_, err = min(vals...)
	if err == nil {
		t.Errorf("err of min(%v) = nil; want error", vals)
	}
}

func TestMax2(t *testing.T) {
	vals := []int{3, 5, -1, -4, 2}
	m := max2(vals[0], vals[1:]...)
	if m != 5 {
		t.Errorf("max(%v) = %d; want %d", vals, m, 5)
	}
}

func TestMin2(t *testing.T) {
	vals := []int{3, 5, -1, -4, 2}
	m := min2(vals[0], vals[1:]...)
	if m != -4 {
		t.Errorf("min(%v) = %d; want %d", vals, m, -4)
	}
}
