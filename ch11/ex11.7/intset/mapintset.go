package intset

import (
	"bytes"
	"fmt"
)

// MapIntSet ...
type MapIntSet struct {
	m map[int]struct{}
}

// NewMapIntSet ...
func NewMapIntSet() *MapIntSet {
	return &MapIntSet{make(map[int]struct{})}
}

// Has ...
func (s *MapIntSet) Has(x int) bool {
	_, ok := s.m[x]
	return ok
}

// Add ...
func (s *MapIntSet) Add(x int) {
	s.m[x] = struct{}{}
}

// UnionWith ...
func (s *MapIntSet) UnionWith(t *MapIntSet) {
	for key := range t.m {
		s.m[key] = struct{}{}
	}
}

// String ...
func (s *MapIntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for key := range s.m {
		if buf.Len() > len("{") {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%d", key)
	}
	buf.WriteByte('}')
	return buf.String()
}
