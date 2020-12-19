package sexpr

import (
	"bytes"
	"testing"
)

func TestDecode(t *testing.T) {
	type iface interface{}

	i := 1
	foo := "foo"
	bar := "bar"
	type testType struct {
		B      bool
		F32    float32
		F64    float64
		C64    complex64
		C128   complex128
		Ptr    *int
		Ptr2   *[]int
		IFace1 interface{}
		IFace2 interface{}
	}
	v := &testType{
		B:      true,
		F32:    123,
		F64:    1.234e+10,
		C64:    1 + 2i,
		C128:   1.2e10 + 1.0e8i,
		Ptr:    &i,
		Ptr2:   &[]int{1, 2},
		IFace1: []int{1, 2, 3},
		IFace2: []map[*string]*[2]int{
			map[*string]*[2]int{&foo: {1, 2}, &bar: {3, 4}},
		},
	}

	b, err := Marshal(v)
	if err != nil {
		t.Errorf("Marshal(%v) failed: %s", v, err)
	}

	reader := bytes.NewReader(b)
	decoder := NewDecoder(reader)

	var decoded testType
	if err := decoder.Decode(&decoded); err != nil {
		t.Fatalf("failed to decode %s: %v", b, err)
	}

	db, err := Marshal(decoded)
	t.Logf("%s\n", db)

	if string(db) != string(b) {
		t.Errorf("Decode() got %s, want %s", db, b)
	}

	//if !reflect.DeepEqual(&decoded, v) {
	//	t.Errorf("Decode() got %v, want %v", decoded, v)
	//}
}
