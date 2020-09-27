package tempconv

import (
	"flag"
	"fmt"
)

// Celsius ...
type Celsius float64

// Fahrenheit ...
type Fahrenheit float64

// Kelvin ...
type Kelvin float64

// CToF ...
func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9.0/5.0 + 32.0) }

// FToC ...
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32.0) * 5.0 / 9.0) }

// CToK ...
func CToK(c Celsius) Kelvin { return Kelvin(c + 273.15) }

// KToC ...
func KToC(k Kelvin) Celsius { return Celsius(k - 273.15) }

func (c Celsius) String() string { return fmt.Sprintf("%g°C", c) }
func (k Kelvin) String() string  { return fmt.Sprintf("%gK", k) }

type celsiusFlag struct{ Celsius }

func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit)
	switch unit {
	case "C", "°C":
		f.Celsius = Celsius(value)
		return nil
	case "F", "°F":
		f.Celsius = FToC(Fahrenheit(value))
		return nil
	case "K":
		f.Celsius = KToC(Kelvin(value))
		return nil
	}
	return fmt.Errorf("invalid temperature %q", s)
}

// CelsiusFlag ...
func CelsiusFlag(name string, value Celsius, usage string) *Celsius {
	f := celsiusFlag{value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Celsius
}
