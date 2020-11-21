package zip

import (
	"archive/zip"
	"bytes"
	"gopl/ch10/ex10.2/extract"
	"io"
	"io/ioutil"
)

type reader struct {
	zr        *zip.Reader
	currIndex int
	name      string
	buf       []byte
	bufIndex  int
}

func init() {
	extract.RegisterFormat("zip", "PK", 0, func(r io.Reader) (extract.Archive, error) {
		buf, err := ioutil.ReadAll(r)
		readerAt := bytes.NewReader(buf)
		zr, err := zip.NewReader(readerAt, int64(len(buf)))

		if err != nil {
			return nil, err
		}
		return &reader{zr: zr, currIndex: -1}, nil
	})
}

func (r *reader) Next() error {
	l := len(r.zr.File)
	if r.currIndex >= l-1 {
		return io.EOF
	}

	r.currIndex++

	currFile := r.zr.File[r.currIndex]
	r.name = currFile.Name
	f, err := currFile.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	r.bufIndex = 0
	r.buf = buf
	return nil
}

func (r *reader) Read(b []byte) (int, error) {
	if r.bufIndex >= len(r.buf) {
		return 0, io.EOF
	}

	n := copy(b, r.buf[r.bufIndex:])
	r.bufIndex += n
	return n, nil
}

func (r *reader) Name() string {
	return r.name
}
