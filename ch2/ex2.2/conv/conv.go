package conv

import "fmt"

type Meter float64
type Feet float64

func (m Meter) String() string {
	return fmt.Sprintf("%gM", m)
}

func (f Feet) String() string {
	return fmt.Sprintf("%gF", f)
}

func (m Meter) ToFeet() Feet {
	return Feet(m * 3.28084)
}

func (f Feet) ToMeter() Meter {
	return Meter(f / 3.28084)
}
