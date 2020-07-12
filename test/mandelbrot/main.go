package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/cmplx"
	"os"
)

const (
	// number of iterations
	maxEsc = 100
	// size of complex plane
	rMin = -2
	rMax = 2.0
	iMin = -2
	iMax = 2
	// width of image
	width = 750
	// color (constant for now)
	red   = 230
	green = 235
	blue  = 255
)

func mandelbrot(c complex128) float64 {
	i := 0
	// iterate the point with f(z) = z^2 + c
	for z := c; cmplx.Abs(z) < 2 && i < maxEsc; i++ {
		z = cmplx.Pow(z, complex(-2, -1.234)) + c
	}
	return float64(maxEsc-i) / maxEsc // how many iterations out of the max did it take to converge
}

func getColor(x float64) color.NRGBA {
	return color.NRGBA{
		uint8(red * x),
		uint8(green * x),
		uint8(blue * x),
		255,
	}
}

func main() {
	// define the image
	scale := width / (rMax - rMin)
	height := int(scale * (iMax - iMin))
	bounds := image.Rect(0, 0, width, height)

	img := image.NewNRGBA(bounds)

	draw.Draw(img, bounds, image.NewUniform(color.Black), image.ZP, draw.Src) // not too sure how this one works

	// fill in the plot, mapping the coordinates of each pixel in the image to
	// a point on the complex plane as defined by our bounds. This is a simple linear transformation
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// map pixel location to point on complex plane.
			pt := mandelbrot(
				complex(
					float64(x)/scale+rMin,
					float64(y)/scale+iMin,
				),
			)

			img.Set(x, y, getColor(pt))
		}
	}

	f, err := os.Create("Mandelbrot.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	if err = png.Encode(f, img); err != nil {
		fmt.Println(err)
		return
	}
}
