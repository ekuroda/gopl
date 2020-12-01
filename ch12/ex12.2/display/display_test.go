package display

import (
	"bytes"
	"strings"
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

func TestMaxDepth(t *testing.T) {
	o := out
	md := maxDepth
	defer func() {
		out = o
		maxDepth = md
	}()

	maxDepth = 5
	b := bytes.NewBuffer(nil)
	out = b
	a := make([]interface{}, 0)
	a = append(a, &a)
	Display("a", a)

	lines := strings.Split(b.String(), "\n")
	if len(lines) != 5 {
		t.Fatalf("len(%q) = %d, want %d", lines, len(lines), 4)
	}

	third := "(*a[0].value)[0].type = *[]interface {}"
	forth := "(*a[0].value)[0].value = *[]interface {}0x"
	if lines[2] != third {
		t.Errorf("third line is %q, want %q", lines[2], third)
	}
	if !strings.HasPrefix(lines[3], forth) {
		t.Errorf("forth line %q shoud has prefix %q", lines[3], forth)
	}
}
