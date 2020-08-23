package main

import (
	"bufio"
	"errors"
	"fmt"
	"gopl/ch2/ex2.2/conv"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	al := len(os.Args)
	switch al {
	case 1:
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			numAndUnit := strings.Split(s.Text(), " ")
			if len(numAndUnit) != 2 {
				log.Fatalf("each line must be number and unit separated by whitespace.")
			}
			converted := convert(numAndUnit[0], numAndUnit[1])
			fmt.Printf("%s %s\t%s\n", numAndUnit[0], numAndUnit[1], converted)
		}
	case 3:
		converted := convert(os.Args[1], os.Args[2])
		fmt.Printf("%s %s\t%s\n", os.Args[1], os.Args[2], converted)
	default:
		log.Fatalf("invalid arguments")
	}

}

func convert(numberString, unit string) string {
	number, err := strconv.ParseFloat(numberString, 64)
	if err != nil {
		log.Fatalf("first part of line must be a number.")
	}
	converted, err := convertNumber(number, unit)
	if err != nil {
		log.Fatalf("fail to convert: number=%f unit=%s, err=%s", number, unit, err)
	}
	return converted
}

func convertNumber(number float64, unit string) (string, error) {
	switch unit {
	case "M", "m":
		f := conv.Meter(number).ToFeet()
		return f.String(), nil
	case "F", "f":
		m := conv.Feet(number).ToMeter()
		return m.String(), nil
	default:
		return "", errors.New("unit must be M/m or F/f.")
	}
}
