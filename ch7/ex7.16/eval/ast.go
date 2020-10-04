package eval

// Expr ...
type Expr interface {
	Eval(env Env) float64
	Check(vars map[Var]bool) error
	String() string
	CollectVars([]Var) []Var
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
