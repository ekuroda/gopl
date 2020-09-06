package main

import "testing"

func TestStringsJoin(t *testing.T) {
	tests := []struct {
		vals []string
		sep  string
		want string
	}{
		{
			vals: []string{"aa", "bb", "cc"},
			sep:  "**",
			want: "aa**bb**cc",
		},
		{
			vals: []string{"aa", "bb", "cc"},
			sep:  "",
			want: "aabbcc",
		},
		{
			vals: []string{},
			sep:  ",",
			want: "",
		},
	}

	for _, test := range tests {
		if s := stringsJoin(test.sep, test.vals...); s != test.want {
			t.Errorf("stringsJoin(%q, %v) = %q; want %q", test.sep, test.vals, s, test.want)
		}
	}
}
