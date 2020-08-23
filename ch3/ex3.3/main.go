package main

import (
	"fmt"
	"math"
)

const (
	width, height = 600, 320
	cells         = 100
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

type polygon struct {
	ax float64
	ay float64
	bx float64
	by float64
	cx float64
	cy float64
	dx float64
	dy float64
	z  float64
}

func main() {
	polygons := make([]polygon, 0)
	zMin, zMax := math.MaxFloat64, -math.MaxFloat32
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>\n", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, az, aok := corner(i+1, j)
			bx, by, bz, bok := corner(i, j)
			cx, cy, cz, cok := corner(i, j+1)
			dx, dy, dz, dok := corner(i+1, j+1)
			if !aok || !bok || !cok || !dok {
				continue
			}
			z := (az + bz + cz + dz) / 4
			if z < zMin {
				zMin = z
			}
			if z > zMax {
				zMax = z
			}
			polygons = append(polygons, polygon{
				ax: ax,
				ay: ay,
				bx: bx,
				by: by,
				cx: cx,
				cy: cy,
				dx: dx,
				dy: dy,
				z:  z,
			})
		}
	}

	for _, p := range polygons {
		// こうじゃなく極大値で赤255/極小値で青255になるようにするということだと思うが
		red := 255 * (p.z - zMin) / (zMax - zMin)
		blue := 255 - red
		color := fmt.Sprintf("#%02x00%02x", int(red), int(blue))
		fmt.Printf("<polygon style='fill: %s' points='%g,%g %g,%g %g,%g %g,%g'/>\n",
			color, p.ax, p.ay, p.bx, p.by, p.cx, p.cy, p.dx, p.dy)
	}

	fmt.Printf("</svg>")
}

func corner(i, j int) (float64, float64, float64, bool) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	z := f(x, y)
	if math.IsNaN(z) || math.IsInf(z, 0) || math.IsInf(z, -1) {
		return 0, 0, 0, false
	}

	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale

	return sx, sy, z, true
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}
