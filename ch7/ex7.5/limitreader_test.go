package limitreader

import (
	"io"
	"strings"
	"testing"
)

func TestLimitReader(t *testing.T) {
	tests := []struct {
		text  string
		limit int
		isEOF bool
		want  int
	}{
		{"0123456789", 3, true, 3},
		{"0123", 4, true, 4},
		{"0", 3, false, 1},
	}

	for _, test := range tests {
		r := LimitReader(strings.NewReader(test.text), test.limit)
		bytes := make([]byte, len([]byte(test.text)))
		n, err := r.Read(bytes)
		if test.isEOF {
			if err != io.EOF {
				t.Fatalf("LimitReader(%q, %d); err = %v, want %v", test.text, test.limit, err, io.EOF)
			}
		} else {
			if err != nil {
				t.Fatalf("failed to read: %v", err)
			}
		}

		if n != test.want {
			t.Errorf("LimitReader(%q, %d); r.Read() = %d, want %d", test.text, test.limit, n, test.want)
		}
	}
}
