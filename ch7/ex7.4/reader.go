package main

import (
	"io"
)

type reader struct {
	s string
}

func NewReader(s string) io.Reader {
	return &reader{s}
}

func (r *reader) Read(p []byte) (n int, err error) {
	end := len(r.s)
	if cap(p) < end {
		end = cap(p)
	}
	n = copy(p, r.s[:end])
	//fmt.Printf("cap=%d, len=%d, end=%d, n=%d\n", cap(p), len(r.s), end, n)
	r.s = r.s[n:]
	if len(r.s) == 0 {
		err = io.EOF
	}
	return
}
