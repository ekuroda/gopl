package main

import (
	"bufio"
	"fmt"
	"gopl/ch7/ex7.15/eval"
	"log"
	"os"
)

func main() {
	var sc = bufio.NewScanner(os.Stdin)

	fmt.Printf("input expr:\n")
	if !sc.Scan() {
		log.Printf("failed to read expr")
		os.Exit(1)
	}
	exprInput := sc.Text()
	expr, err := eval.Parse(exprInput)
	if err != nil {
		log.Printf("failed to parse expr %s: %s", exprInput, err)
		os.Exit(1)
	}

	check := make(map[eval.Var]bool)
	if err = expr.Check(check); err != nil {
		log.Printf("failed to check expr %s: %s", expr, err)
		os.Exit(1)
	}

	vars := make(map[eval.Var]float64)
	if err = expr.InputVars(vars); err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}

	r := expr.Eval(eval.Env(vars))
	fmt.Printf("=> %g\n", r)
}
