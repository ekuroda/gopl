package cycle

import (
	"reflect"
	"unsafe"
)

func hasCycle(v reflect.Value, seen map[element]bool) bool {
	if !v.IsValid() {
		return false
	}

	if v.CanAddr() {
		vptr := unsafe.Pointer(v.UnsafeAddr())
		c := element{vptr, v.Type()}
		if seen[c] {
			return true
		}
		seen[c] = true
	}

	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return hasCycle(v.Elem(), seen)
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if hasCycle(v.Index(i), seen) {
				return true
			}
		}
		return false
	case reflect.Struct:
		for i, n := 0, v.NumField(); i < n; i++ {
			if hasCycle(v.Field(i), seen) {
				return true
			}
		}
		return false
	case reflect.Map:
		for _, k := range v.MapKeys() {
			if hasCycle(k, seen) || hasCycle(v.MapIndex(k), seen) {
				return true
			}
		}
		return false
	}
	return false
}

// HasCycle ...
func HasCycle(v interface{}) bool {
	seen := make(map[element]bool)
	return hasCycle(reflect.ValueOf(v), seen)
}

type element struct {
	v unsafe.Pointer
	// unsafe.Pointerのアドレス値だけの比較では不足
	// ある型の変数のアドレスと、その変数の先頭の要素のアドレスが同じになる
	t reflect.Type
}
