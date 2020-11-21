package tar

import (
	"archive/tar"
	"gopl/ch10/ex10.2/extract"
	"io"
)

type reader struct {
	tr   *tar.Reader
	name string
}

func init() {
	extract.RegisterFormat("tar", "ustar", 257, func(r io.Reader) (extract.Archive, error) {
		tr := tar.NewReader(r)
		return &reader{tr: tr}, nil
	})
}

func (r *reader) Next() error {
	h, err := r.tr.Next()
	if err != nil {
		return err
	}
	r.name = h.Name
	return nil
}

func (r *reader) Read(b []byte) (int, error) {
	return r.tr.Read(b)
}

func (r *reader) Name() string {
	return r.name
}
