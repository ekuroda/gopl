package main

import "fmt"

func main() {
	s := panicRecover()
	fmt.Println(s)
}

func panicRecover() (s string) {
	defer func() {
		recover()
		s = "recover"
	}()
	panic("fatal")
}
