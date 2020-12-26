package numequal

import (
	"math"
	"reflect"
)

const epsilon = 1e-9

func equal(x, y reflect.Value) bool {
	if !x.IsValid() || !y.IsValid() {
		return x.IsValid() == y.IsValid()
	}
	if x.Type() != y.Type() {
		return false
	}

	switch x.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:
		return x.Int() == y.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		return x.Uint() == y.Uint()
	case reflect.Float32, reflect.Float64:
		return math.Abs(x.Float()-y.Float()) < epsilon
	case reflect.Complex64, reflect.Complex128:
		xc, yc := x.Complex(), y.Complex()
		return math.Abs(real(xc)-real(yc)) < epsilon && math.Abs(imag(xc)-imag(yc)) < epsilon
	}
	panic("unreachable")
}

// Equal ...
func Equal(x, y interface{}) bool {
	return equal(reflect.ValueOf(x), reflect.ValueOf(y))
}
