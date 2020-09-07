package intset

import "testing"

func TestAdd(t *testing.T) {
	var x, y IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	want := "{1 9 144}"
	if s := x.String(); s != want {
		t.Errorf("x = %q; want %q", s, want)
	}

	y.Add(9)
	y.Add(42)

	x.UnionWith(&y)
	want = "{1 9 42 144}"
	if s := x.String(); s != want {
		t.Errorf("x = %q; want %q", s, want)
	}

	if b := x.Has(9); b != true {
		t.Errorf("x.Has(9) = %t; want %t", b, true)
	}

	if b := x.Has(123); b != false {
		t.Errorf("x.Has(123) = %t; want %t", b, false)
	}
}

func TestAddAll(t *testing.T) {
	var x IntSet
	x.AddAll(1, 144, 9)
	want := "{1 9 144}"
	if s := x.String(); s != want {
		t.Errorf("x = %q; want %q", s, want)
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{
			nums: []int{},
			want: 0,
		},
		{
			nums: []int{0},
			want: 1,
		},
		{
			nums: []int{2},
			want: 1,
		},
		{
			nums: []int{2, 8},
			want: 2,
		},
		{
			nums: []int{100, 200, 1},
			want: 3,
		},
	}

	for _, test := range tests {
		var x IntSet
		for _, num := range test.nums {
			x.Add(num)
		}
		if n := x.Len(); n != test.want {
			t.Errorf("x.Len() = %d; want %d", n, test.want)
		}
	}
}

func TestRemove(t *testing.T) {
	var x IntSet
	x.Add(100)
	x.Add(2)
	x.Add(3)

	x.Remove(2)
	if x.Has(2) || !x.Has(3) || !x.Has(100) {
		t.Errorf("x.Remove(%d); x.Has(%d), x.Has(%d), x.Has(%d) got (%t, %t, %t); want (%t, %t, %t)",
			2, 2, 3, 100, false, true, true, x.Has(2), x.Has(3), x.Has(100))
	}

	x.Remove(200)
	if x.Has(200) || !x.Has(3) || !x.Has(100) {
		t.Errorf("x.Remove(%d); x.Has(%d), x.Has(%d), x.Has(%d) got (%t, %t, %t); want (%t, %t, %t)",
			200, 200, 3, 100, false, true, true, x.Has(200), x.Has(3), x.Has(100))
	}

	x.Remove(4)
	if x.Has(4) || !x.Has(3) || !x.Has(100) {
		t.Errorf("x.Remove(%d); x.Has(%d), x.Has(%d), x.Has(%d) got (%t, %t, %t); want (%t, %t, %t)",
			4, 4, 3, 100, false, true, true, x.Has(4), x.Has(3), x.Has(100))
	}

	x.Remove(100)
	if x.Has(100) || !x.Has(3) {
		t.Errorf("x.Remove(%d); x.Has(%d), x.Has(%d) got (%t, %t); want (%t, %t)",
			100, 100, 3, false, true, x.Has(100), x.Has(3))
	}
}

func TestClear(t *testing.T) {
	var x IntSet
	x.Add(100)
	x.Add(2)
	x.Add(3)

	x.Clear()
	if x.Len() != 0 {
		t.Errorf("x.Clear(); x.Len() = %d; want %d", x.Len(), 0)
	}
}

func TestCopy(t *testing.T) {
	var x IntSet
	x.Add(100)
	x.Add(2)
	x.Add(3)

	y := x.Copy()

	for i, w := range x.words {
		if w != y.words[i] {
			t.Errorf("y.words[%d] = %d; want %d", i, y.words[i], w)
		}
	}
	if len(y.words) != len(x.words) {
		t.Errorf("len(y.words) = %d; want %d", len(y.words), len(x.words))
	}
}
