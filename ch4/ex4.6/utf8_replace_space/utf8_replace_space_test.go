package utf8_replace_space

import "testing"

func TestReplaceSpace(t *testing.T) {
	s := "\t \tthis\nイズ\t\na\tペン \n"
	copy := s
	want := " this イズ a ペン "
	replaced := ReplaceSpace([]byte(s))

	if string(replaced) != want {
		t.Errorf("ReplaceSpace(%q) = %q; want %q", copy, string(replaced), want)
	}
}
