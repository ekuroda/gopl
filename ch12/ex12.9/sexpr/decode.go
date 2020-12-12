package sexpr

import (
	"fmt"
	"io"
	"strconv"
	"text/scanner"
)

// Token ...
type Token interface{}

// Symbol ...
type Symbol string

// String ...
type String string

// Int ...
type Int int

// StartList ...
type StartList struct{}

// EndList ...
type EndList struct{}

// Decoder ...
type Decoder struct {
	lex *lexer
}

// NewDecoder ...
func NewDecoder(reader io.Reader) *Decoder {
	lex := &lexer{scan: scanner.Scanner{Mode: scanner.GoTokens}}
	lex.scan.Init(reader)
	return &Decoder{lex}
}

// Token ...
func (d *Decoder) Token() (Token, error) {
	var err error
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("error at %s: %v", d.lex.scan.Position, x)
		}
	}()

	t := token(d.lex)
	if d.lex.token == scanner.EOF {
		err = io.EOF
	}
	return t, err
}

type lexer struct {
	scan  scanner.Scanner
	token rune
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }

func (lex *lexer) consume(want rune) {
	if lex.token != want {
		panic(fmt.Sprintf("got %q, want %q", lex.text(), want))
	}
	lex.next()
}

func token(lex *lexer) Token {
	lex.next()
	switch lex.token {
	case scanner.Ident:
		token := Symbol(lex.text())
		return token
	case scanner.String:
		s, _ := strconv.Unquote(lex.text())
		token := String(s)
		return token
	case scanner.Int:
		i, _ := strconv.Atoi(lex.text())
		token := Int(i)
		return token
	case '(':
		return StartList{}
	case ')':
		return EndList{}
	case scanner.EOF:
		return nil
	}
	panic(fmt.Sprintf("unexpected token %q", lex.text()))
}
