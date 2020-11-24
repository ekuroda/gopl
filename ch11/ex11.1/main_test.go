package main

import (
	"bytes"
	"testing"
)

func TestFoo(t *testing.T) {
	tests := []struct {
		input  []byte
		output string
	}{
		{
			input: []byte("あaういéab\nう𠀋cá"),
			output: `rune	count
'\n'	1
'a'	2
'b'	1
'c'	1
'á'	1
'é'	1
'あ'	1
'い'	1
'う'	2
'𠀋'	1

len	count
1	5
2	2
3	4
4	1
`,
		},
		{
			input: []byte("Foo\xe4\xb8"),
			output: `rune	count
'F'	1
'o'	2

len	count
1	3
2	0
3	0
4	0

2 invalid UTF-8 characters
`,
		},
		{
			input: []byte(""),
			output: `rune	count

len	count
1	0
2	0
3	0
4	0
`,
		},
	}

	for _, test := range tests {
		var buf bytes.Buffer
		count(bytes.NewReader(test.input), &buf)
		result := buf.String()
		//fmt.Printf("%q\n", result)
		//fmt.Printf("%q\n", test.output)
		if result != test.output {
			t.Errorf("count(%q) = %q, want %q", test.input, result, test.output)
		}
	}
}
