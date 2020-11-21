package extract

import (
	"bufio"
	"errors"
	"io"
	"sync"
	"sync/atomic"
)

// Archive ...
type Archive interface {
	Next() error
	Read(b []byte) (int, error)
	Name() string
}

// ErrFormat indicates that decoding encountered an unknown format.
var ErrFormat = errors.New("extract: unknown format")

// A format holds an archive format's name, magic header and how to decode it.
type format struct {
	name, magic string
	magicOffset int
	decode      func(io.Reader) (Archive, error)
}

// Formats is the list of registered formats.
var (
	formatsMu     sync.Mutex
	atomicFormats atomic.Value
)

// RegisterFormat registers an archive format for use by Decode.
// Name is the name of the format, like "zip" or "tar".
// Decode is the function that decodes the encoded archive.
func RegisterFormat(name, magic string, magicOffset int, decode func(io.Reader) (Archive, error)) {
	formatsMu.Lock()
	formats, _ := atomicFormats.Load().([]format)
	atomicFormats.Store(append(formats, format{name, magic, magicOffset, decode}))
	formatsMu.Unlock()
}

// A reader is an io.Reader that can also peek ahead.
type reader interface {
	io.Reader
	Peek(int) ([]byte, error)
}

// asReader converts an io.Reader to a reader.
func asReader(r io.Reader) reader {
	if rr, ok := r.(reader); ok {
		return rr
	}
	return bufio.NewReader(r)
}

// Sniff determines the format of r's data.
func sniff(r reader) format {
	formats, _ := atomicFormats.Load().([]format)
	for _, f := range formats {
		p, err := r.Peek(f.magicOffset + len(f.magic))
		if err != nil {
			continue
		}
		if string(p[f.magicOffset:]) == f.magic {
			return f
		}
	}
	return format{}
}

// Decode decodes an archive that has been encoded in a registered format.
// The string returned is the format name used during format registration.
// Format registration is typically done by an init function in the codec-
// specific package.
func Decode(r io.Reader) (Archive, string, error) {
	rr := asReader(r)
	f := sniff(rr)
	if f.decode == nil {
		return nil, "", ErrFormat
	}
	archive, err := f.decode(rr)
	return archive, f.name, err
}
