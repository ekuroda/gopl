package display

import (
	"bytes"
	"testing"
)

func TestMapKey(t *testing.T) {
	o := out
	defer func() { out = o }()

	b := bytes.NewBuffer(nil)
	out = b
	s := map[struct {
		x string
		y int
	}]int{
		{"aaa", 1}: 10,
		{"bbb", 2}: 100,
	}
	Display("s", s)

	sgot := b.String()
	swant := `Display s (map[struct { x string; y int }]int):
s[{x: "aaa", y: 1}] = 10
s[{x: "bbb", y: 2}] = 100`

	if sgot != sgot {
		t.Errorf("Display(%v) = %q, want %q", s, sgot, swant)
	}

	b = bytes.NewBuffer(nil)
	out = b
	a := map[[2]string]int{
		{"foo", "bar"}:   10,
		{"hoge", "fuga"}: 100,
	}
	Display("a", a)

	agot := b.String()
	awant := `Display a (map[[2]string]int):
a[{"foo", "bar"}] = 10
a[{"hoge", "fuga"}] = 100`

	if agot != agot {
		t.Errorf("Display(%v) = %q, want %q", s, agot, awant)
	}
}
