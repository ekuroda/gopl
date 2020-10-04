package eval

// Expr ...
type Expr interface {
	Eval(env Env) float64
	Check(vars map[Var]bool) error
	String() string
	InputVars(map[Var]float64) error
}

// Var ...
type Var string

type literal float64

type unary struct {
	op rune // one of '+', '-'
	x  Expr
}

type binary struct {
	op   rune // one of '+', '-', '*', '/'
	x, y Expr
}

type call struct {
	fn   string // one of "pow", "sin", "sqrt"
	args []Expr
}
