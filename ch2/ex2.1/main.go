package main

import (
	"fmt"
	"gopl/ch2/ex2.1/tempconv"
)

func main() {
	absoluteZeroCInF := tempconv.CToF(tempconv.AbsoluteZeroC)
	fmt.Printf("absoluteZeroCInF:%s\n", absoluteZeroCInF)
	fmt.Printf("absoluteZeroC:%s\n", tempconv.FToC(absoluteZeroCInF))

	absoluteZeroCInK := tempconv.CToK(tempconv.AbsoluteZeroC)
	fmt.Printf("absoluteZeroCInK:%s\n", absoluteZeroCInK)
	fmt.Printf("absoluteZeroC:%s\n", tempconv.KToC(absoluteZeroCInK))
}
