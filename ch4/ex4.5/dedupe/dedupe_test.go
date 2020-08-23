package dedupe

import "testing"

func TestDedupe(t *testing.T) {
	tests := []struct {
		input []string
		want  []string
	}{
		{
			input: []string{},
			want:  []string{},
		},
		{
			input: []string{"aaa"},
			want:  []string{"aaa"},
		},
		{
			input: []string{"aaa", "bbb", "bbb", "bbb", "ccc", "aaa", "aaa"},
			want:  []string{"aaa", "bbb", "ccc", "aaa"},
		},
	}

	for _, test := range tests {
		got := Dedupe(test.input)
		if !equal(got, test.want) {
			t.Errorf("Dedupe(%v) = %v; want %v", test.input, got, test.want)
		}
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, s := range a {
		if s != b[i] {
			return false
		}
	}

	return true
}
