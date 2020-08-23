package reverse_utf8

import (
	"testing"
)

func TestReverseUtf8(t *testing.T) {
	s := "\t \tthis\nイズ\t\na\tペン \n"
	copy := s
	want := "\n ンペ\ta\n\tズイ\nsiht\t \t"
	b := []byte(s)
	ReverseUTF8(b)

	if string(b) != want {
		t.Errorf("ReverseUTF8(%q) = %q; want %q", copy, string(b), want)
	}
}
