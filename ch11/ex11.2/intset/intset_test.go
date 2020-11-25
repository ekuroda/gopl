package intset

import (
	"testing"
)

func TestAdd(t *testing.T) {
	var x IntSet
	mx := NewMapIntSet()

	elems := []int{1, 2, 5}
	for _, elem := range elems {
		x.Add(elem)
		mx.Add(elem)
	}

	elems = []int{1, 2, 3, 4, 5, 6}
	for _, elem := range elems {
		xhas := x.Has(elem)
		mxhas := mx.Has(elem)
		if xhas != mxhas {
			t.Errorf("after Add %v; x.Has(%d) = %t, mx.Has(%d) = %t; not equal", elems, elem, xhas, elem, mxhas)
		}
	}

	var y IntSet
	my := NewMapIntSet()
	joinElems := []int{3, 6}
	for _, elem := range joinElems {
		y.Add(elem)
		my.Add(elem)
	}

	x.UnionWith(&y)
	mx.UnionWith(my)
	for _, elem := range elems {
		xhas := x.Has(elem)
		mxhas := mx.Has(elem)
		if xhas != mxhas {
			t.Errorf("after UnionWith %v; x.Has(%d) = %t, mx.Has(%d) = %t; not equal", joinElems, elem, xhas, elem, mxhas)
		}
	}
}
