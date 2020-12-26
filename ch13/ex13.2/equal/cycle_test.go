package cycle

import (
	"bytes"
	"testing"
)

func TestEqual(t *testing.T) {
	one := 1

	type CyclePtr *CyclePtr
	var cyclePtr CyclePtr
	cyclePtr = &cyclePtr

	type CycleSlice []CycleSlice
	var cycleSlice = make(CycleSlice, 1)
	cycleSlice[0] = cycleSlice

	ch := make(chan int)

	type mystring string

	var iface interface{} = &one

	type St1 struct{}
	type St2 struct{ v *St2 }
	var st St2
	st.v = &st

	//type Map0 map[string]Map0
	//m0 := make(map[string]Map0)
	//m0[""] = m0

	type Map1 map[string]CyclePtr
	m1 := map[string]CyclePtr{"": cyclePtr}

	type Map2 map[CyclePtr]int
	m2 := map[CyclePtr]int{cyclePtr: 0}

	for _, test := range []struct {
		v    interface{}
		want bool
	}{
		{1, false},
		{"foo", false},
		{mystring("foo"), false},
		{[]string{"foo"}, false},
		{[]string{}, false},
		{[]string(nil), false},
		{cycleSlice, true},
		{map[string][]int{"foo": {1, 2, 3}}, false},
		{map[string][]int{}, false},
		{map[string][]int(nil), false},
		//{m0, true}, // stack overflow
		{m1, true},
		{m2, true},
		{&one, false},
		{new(bytes.Buffer), false},
		{cyclePtr, true},
		{(func())(nil), false},
		{func() {}, false},
		{[...]int{1, 2, 3}, false},
		{ch, false},
		{&iface, false},
		{St1{}, false},
		{St2{}, false},
		{st, true},
	} {
		if HasCycle(test.v) != test.want {
			t.Errorf("HasCycle(%#v) = %t", test.v, !test.want)
		}
	}
}
