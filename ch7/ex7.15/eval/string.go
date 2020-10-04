package eval

import (
	"fmt"
)

// String ...
func (v Var) String() string {
	return string(v)
}

// String ...
func (l literal) String() string {
	return fmt.Sprintf("%g", float64(l))
}

// String ...
func (u unary) String() string {
	return fmt.Sprintf("(%c %s)", u.op, u.x)
}

// String ...
func (b binary) String() string {
	return fmt.Sprintf("(%s %c %s)", b.x, b.op, b.y)
}

// String ...
func (c call) String() string {
	switch c.fn {
	case "pow":
		return fmt.Sprintf("pow(%s, %s)", c.args[0], c.args[1])
	case "sin":
		return fmt.Sprintf("sin(%s)", c.args[0])
	case "sqrt":
		return fmt.Sprintf("sqrt(%s)", c.args[0])
	}
	panic(fmt.Sprintf("unsupported function call: %s", c.fn))
}
