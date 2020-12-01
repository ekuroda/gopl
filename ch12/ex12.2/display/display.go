package display

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var out io.Writer = os.Stdout
var maxDepth = 5

// Display ...
func Display(name string, x interface{}) {
	fmt.Fprintf(out, "Display %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x), 0)
}

func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr,
		reflect.Slice, reflect.Map:
		return v.Type().String() + "0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	default:
		return v.Type().String() + " value"
	}
}

func formatMapKey(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Struct:
		var b strings.Builder
		b.WriteString("{")
		for i := 0; i < v.NumField(); i++ {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(fmt.Sprintf("%s: %s",
				v.Type().Field(i).Name, formatAtom(v.Field(i))))
		}
		b.WriteString("}")
		return b.String()
	case reflect.Array:
		var b strings.Builder
		b.WriteString("{")
		for i := 0; i < v.Len(); i++ {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(formatAtom(v.Index(i)))
		}
		b.WriteString("}")
		return b.String()
	default:
		return formatAtom(v)
	}
}

func display(path string, v reflect.Value, depth int) {
	if depth >= maxDepth {
		fmt.Fprintf(out, "%s = %s\n", path, formatAtom(v))
		return
	}

	switch v.Kind() {
	case reflect.Invalid:
		fmt.Fprintf(out, "%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i), depth+1)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(fieldPath, v.Field(i), depth+1)
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path,
				formatMapKey(key)), v.MapIndex(key), depth+1)
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Fprintf(out, "%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem(), depth+1)
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Fprintf(out, "%s = nil\n", path)
		} else {
			fmt.Fprintf(out, "%s.type = %s\n", path, v.Elem().Type())
			display(path+".value", v.Elem(), depth+1)
		}
	default:
		fmt.Fprintf(out, "%s = %s\n", path, formatAtom(v))
	}
}
