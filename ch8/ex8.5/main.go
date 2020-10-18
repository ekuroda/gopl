package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	img := createImage(runtime.NumCPU())
	png.Encode(os.Stdout, img)
}

type p struct {
	px int
	py int
	y  float64
}

func createImage(workerNum int) *image.RGBA {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
	)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	pyc := make(chan int, workerNum*2)
	go func() {
		for py := 0; py < height; py++ {
			pyc <- py
		}
		close(pyc)
	}()

	var wg sync.WaitGroup
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for py := range pyc {
				y := float64(py)/height*(ymax-ymin) + ymin
				for px := 0; px < width; px++ {
					x := float64(px)/width*(xmax-xmin) + xmin
					z := complex(x, y)
					img.Set(px, py, mandelbrot(z))
				}
			}
		}()
	}

	wg.Wait()
	return img
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			c := uint64(math.Pow(2, 24) * math.Log(float64(n)) / math.Log(float64(iterations)))
			r := c & 255
			c >>= 8
			g := c & 255
			c >>= 8
			b := c & 255
			return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
		}
	}
	return color.Black
}
