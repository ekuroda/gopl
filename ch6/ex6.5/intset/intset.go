package intset

import (
	"bytes"
	"fmt"
)

var uintsize = 32 << (^uint(0) >> 63)

// IntSet ..
type IntSet struct {
	words []uint64
}

// Has ...
func (s *IntSet) Has(x int) bool {
	word, bit := x/uintsize, uint(x%uintsize)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add ...
func (s *IntSet) Add(x int) {
	word, bit := x/uintsize, uint(x%uintsize)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// AddAll ...
func (s *IntSet) AddAll(xs ...int) {
	for _, x := range xs {
		s.Add(x)
	}
}

// UnionWith ...
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// IntersectWith ...
func (s *IntSet) IntersectWith(t *IntSet) {
	for i := range s.words {
		if i < len(t.words) {
			s.words[i] &= t.words[i]
		} else {
			s.words[i] = 0
		}
	}
}

// DifferenceWith ...
func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &= ^tword
		}
	}
}

// SymmetricDifference ...
func (s *IntSet) SymmetricDifference(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			//s.words[i] = (s.words[i] | tword) & ^(s.words[i] & tword)
			s.words[i] ^= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < uintsize; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

// Elems ...
func (s *IntSet) Elems() []int {
	var elems []int
	for i, word := range s.words {
		for j := 0; word != 0; j++ {
			if word&1 > 0 {
				elems = append(elems, i*64+j)
			}
			word >>= 1
		}
	}
	return elems
}

// Len ...
func (s *IntSet) Len() int {
	count := 0
	for _, word := range s.words {
		for word != 0 {
			if word&1 > 0 {
				count++
			}
			word >>= 1
		}
	}
	return count
}

// Remove ...
func (s *IntSet) Remove(x int) {
	word, bit := x/uintsize, uint(x%uintsize)
	if word < len(s.words) {
		s.words[word] &= ^(1 << bit)
	}
}

// Clear ...
func (s *IntSet) Clear() {
	s.words = make([]uint64, 0)
}

// Copy ...
func (s *IntSet) Copy() *IntSet {
	words := make([]uint64, len(s.words))
	copy(words, s.words)

	return &IntSet{
		words: words,
	}
}
