package main

import "testing"

func TestIsAnagram(t *testing.T) {
	var tests = []struct {
		a    string
		b    string
		want bool
	}{
		{"", "", false},
		{"abc", "abc", false},
		{"aac", "bcaa", false},
		{"aacc", "ccac", false},
		{"baacab", "ababca", true},
	}

	for _, test := range tests {
		if got := isAnagram(test.a, test.b); got != test.want {
			t.Errorf("isAnagram(%q, %q) = %t; want %t", test.a, test.b, got, test.want)
		}
	}
}
