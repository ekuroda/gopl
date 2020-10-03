package palindrome

import "testing"

type runeSlice []rune

func (s runeSlice) Len() int           { return len(s) }
func (s runeSlice) Less(i, j int) bool { return s[i] != s[j] }
func (s runeSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"", true},
		{"あ", true},
		{"いい", true},
		{"しんぶんし", true},
		{"あいうあ", false},
		{"あいうおあ", false},
	}

	for _, test := range tests {
		ok := IsPalindrome(runeSlice([]rune(test.s)))
		if ok != test.want {
			t.Errorf("IsPalindrome(%q) = %t, want %t", test.s, ok, test.want)
		}
	}
}
