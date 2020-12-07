package sexpr

import (
	"bytes"
	"fmt"
	"reflect"
)

// Marshal ...
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")
	case reflect.Bool:
		if v.Bool() {
			buf.WriteString("t")
		} else {
			buf.WriteString("nil")
		}
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%g", v.Float())
	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())
	case reflect.Ptr:
		return encode(buf, v.Elem())
	case reflect.Array, reflect.Slice:
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			if err := encode(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(')')
	case reflect.Struct:
		buf.WriteByte('(')
		j := 0
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			if isZero(f) {
				continue
			}
			if j > 0 {
				buf.WriteByte(' ')
			}
			j++
			fmt.Fprintf(buf, "(%s ", v.Type().Field(i).Name)
			if err := encode(buf, f); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	case reflect.Map:
		buf.WriteByte('(')
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteByte('(')
			if err := encode(buf, key); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := encode(buf, v.MapIndex(key)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	default:
		return fmt.Errorf("unsupported type: %s", v.Type())
	}

	return nil
}

func isZero(v reflect.Value) bool {
	zero := reflect.Zero(v.Type())
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool() == zero.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return v.Int() == zero.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == zero.Uint()
	case reflect.Float32, reflect.Float64:
		return v.Float() == zero.Float()
	case reflect.String:
		return v.String() == zero.String()
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !isZero(v.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !isZero(v.Field(i)) {
				return false
			}
		}
		return true
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Ptr:
		return v.IsNil() || isZero(reflect.Indirect(v))
	default:
		return false
	}
}
