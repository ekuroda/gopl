package limitreader

import (
	"io"
)

type reader struct {
	limit  int
	count  int
	reader io.Reader
}

// LimitReader ...
func LimitReader(r io.Reader, n int) io.Reader {
	return &reader{n, 0, r}
}

func (r *reader) Read(p []byte) (n int, err error) {
	end := cap(p)
	if end > r.limit {
		end = r.limit
	}
	n, err = r.reader.Read(p[:end])
	//fmt.Printf("cap=%d, end=%d, n=%d\n", cap(p), end, n)
	if err != nil {
		return
	}

	r.count += n
	if r.count >= r.limit {
		err = io.EOF
	}
	return n, err
}
