package sexpr

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"
)

// Decoder ...
type Decoder struct {
	reader io.Reader
}

// NewDecoder ...
func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{reader}
}

// Decode ...
func (d *Decoder) Decode(out interface{}) (err error) {
	lex := &lexer{scan: scanner.Scanner{Mode: scanner.GoTokens}}
	lex.scan.Init(d.reader)
	lex.next()
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("error at %s: %v", lex.scan.Position, x)
		}
	}()
	read(lex, reflect.ValueOf(out).Elem())
	return nil
}

type lexer struct {
	scan  scanner.Scanner
	token rune
}

func (lex *lexer) next() {
	//fmt.Printf("%s\n", lex.text())
	lex.token = lex.scan.Scan()
}
func (lex *lexer) text() string { return lex.scan.TokenText() }

func (lex *lexer) consume(want rune) {
	if lex.token != want {
		panic(fmt.Sprintf("got %v, want %v", lex.text(), want))
	}
	lex.next()
}

func read(lex *lexer, v reflect.Value) {
	if v.Kind() == reflect.Ptr {
		value := reflect.New(v.Type().Elem())
		v.Set(value)
		v = reflect.Indirect(value)
	}

	switch lex.token {
	case scanner.Ident:
		switch lex.text() {
		case "nil":
			v.Set(reflect.Zero(v.Type()))
		case "t":
			v.SetBool(true)
		}
		lex.next()
		return
	case scanner.String:
		s, _ := strconv.Unquote(lex.text())
		v.SetString(s)
		lex.next()
		return
	case scanner.Int:
		// 負数未対応
		i, _ := strconv.Atoi(lex.text())
		switch v.Type().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v.SetInt(int64(i))
		case reflect.Float32, reflect.Float64:
			v.SetFloat(float64(i))
		}
		lex.next()
		return
	case scanner.Float:
		switch v.Kind() {
		case reflect.Float32:
			f, _ := strconv.ParseFloat(lex.text(), 32)
			v.SetFloat(f)
		case reflect.Float64:
			f, _ := strconv.ParseFloat(lex.text(), 64)
			v.SetFloat(f)
		}
		lex.next()
		return
	case '#': // #C(1.0, 2.0)
		lex.consume('#')
		lex.consume(scanner.Ident)
		lex.consume('(')
		r, _ := strconv.ParseFloat(lex.text(), 32)
		lex.next()
		i, _ := strconv.ParseFloat(lex.text(), 32)
		lex.next()
		lex.consume(')')
		v.SetComplex(complex(r, i))
		return
	case '(':
		lex.next()
		readList(lex, v)
		lex.next() // consume ')'
		return
	}
	panic(fmt.Sprintf("unexpected token %q", lex.text()))
}

func readList(lex *lexer, v reflect.Value) {
	//fmt.Printf("kind: %s\n", v.Kind())
	switch v.Kind() {
	case reflect.Array:
		for i := 0; !endList(lex); i++ {
			read(lex, v.Index(i))
		}
	case reflect.Slice:
		for !endList(lex) {
			item := reflect.New(v.Type().Elem()).Elem()
			read(lex, item)
			v.Set(reflect.Append(v, item))
		}
	case reflect.Struct:
		for !endList(lex) {
			lex.consume('(')
			if lex.token != scanner.Ident {
				panic(fmt.Sprintf("got token %q, want field name", lex.text()))
			}
			name := lex.text()
			lex.next()
			read(lex, v.FieldByName(name))
			lex.consume(')')
		}
	case reflect.Map:
		v.Set(reflect.MakeMap(v.Type()))
		for !endList(lex) {
			lex.consume('(')
			key := reflect.New(v.Type().Key()).Elem()
			read(lex, key)
			value := reflect.New(v.Type().Elem()).Elem()
			read(lex, value)
			v.SetMapIndex(key, value)
			lex.consume(')')
		}
	case reflect.Interface:
		t, _ := strconv.Unquote(lex.text())
		value := reflect.New(parseType(t)).Elem()
		lex.next()
		read(lex, value)
		v.Set(value)
	default:
		panic(fmt.Sprintf("cannot decode list into %v", v.Type()))
	}
}

func endList(lex *lexer) bool {
	switch lex.token {
	case scanner.EOF:
		panic("end of file")
	case ')':
		return true
	}
	return false
}

var typeMap = map[string]reflect.Type{
	"int":        reflect.TypeOf(int(0)),
	"int8":       reflect.TypeOf(int8(0)),
	"int16":      reflect.TypeOf(int16(0)),
	"int32":      reflect.TypeOf(int32(0)),
	"int64":      reflect.TypeOf(int64(0)),
	"uint":       reflect.TypeOf(uint(0)),
	"uint8":      reflect.TypeOf(uint8(0)),
	"uint16":     reflect.TypeOf(uint16(0)),
	"uint32":     reflect.TypeOf(uint32(0)),
	"uint64":     reflect.TypeOf(uint64(0)),
	"bool":       reflect.TypeOf(false),
	"string":     reflect.TypeOf(""),
	"complex64":  reflect.TypeOf(complex64(0 + 0i)),
	"complex128": reflect.TypeOf(complex128(0 + 0i)),
}

func parseType(s string) reflect.Type {
	t, ok := typeMap[s]
	if ok {
		return t
	}

	if strings.HasPrefix(s, "*") {
		return reflect.PtrTo(parseType(s[1:]))
	}

	if strings.HasPrefix(s, "[]") {
		return reflect.SliceOf(parseType(s[2:]))
	}

	// [[]int]int, map[[1]int]intといった[...]のネストには未対応

	if s[0] == '[' {
		i := strings.Index(s, "]")
		if i > 0 {
			len, _ := strconv.Atoi(s[1:i])
			return reflect.ArrayOf(len, parseType(s[i+1:]))
		}
	}

	if strings.HasPrefix(s, "map") {
		i := strings.Index(s, "]")
		if i > 0 {
			return reflect.MapOf(parseType(s[4:i]), parseType(s[i+1:]))
		}
	}

	panic(fmt.Sprintf("cannot parse type string: %s", s))
}
