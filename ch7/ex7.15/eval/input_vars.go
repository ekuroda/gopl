package eval

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

var sc = bufio.NewScanner(os.Stdin)

// InputVars ...
func (v Var) InputVars(vars map[Var]float64) error {
	if _, ok := vars[v]; ok {
		return nil
	}

	fmt.Printf("input var: %s\n", v)

	if !sc.Scan() {
		log.Printf("failed to read input")
		return fmt.Errorf("failed to read input")
	}

	s := sc.Text()
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Printf("failed to parse input %s: %s", s, err)
		return fmt.Errorf("failed to parse input %s", s)
	}
	vars[v] = val

	return nil
}

// InputVars ...
func (literal) InputVars(vars map[Var]float64) error {
	return nil
}

// InputVars ...
func (u unary) InputVars(vars map[Var]float64) error {
	return u.x.InputVars(vars)
}

// InputVars ...
func (b binary) InputVars(vars map[Var]float64) error {
	if err := b.x.InputVars(vars); err != nil {
		return err
	}
	return b.y.InputVars(vars)
}

// InputVars ...
func (c call) InputVars(vars map[Var]float64) error {
	for _, arg := range c.args {
		if err := arg.InputVars(vars); err != nil {
			return err
		}
	}
	return nil
}
