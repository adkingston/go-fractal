package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math/cmplx"
	"os"
)

const (
	// number of iterations
	maxEsc = 200
	// size of complex plane
	rMin = -2
	rMax = 2.0
	iMin = -2
	iMax = 2
	// width of image
	width = 1750
	// color (constant for now)
	red   = 230
	green = 235
	blue  = 255
)

var (
	palette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0x00, 0x00, 0xff, 0xff},
		color.RGBA{0x00, 0xff, 0x00, 0xff},
		color.RGBA{0x00, 0xff, 0xff, 0xff},
		color.RGBA{0xff, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0x00, 0xff, 0xff},
		color.RGBA{0xff, 0xff, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
	}
)

func mandelbrot(c complex128, n float64) float64 {
	i := 0
	// iterate the point with f(z) = z^2 + c
	for z := c; cmplx.Abs(z) < 2 && i < maxEsc; i++ {
		z = cmplx.Pow(z, complex(n+1, 0)) + complex(n, 0)*cmplx.Pow(z, complex(2, 0)) + c
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

	var images []*image.Paletted
	var delays []int

	// fill in the plot, mapping the coordinates of each pixel in the image to
	// a point on the complex plane as defined by our bounds. This is a simple linear transformation
	for n := 0.0; n <= 10.0; n += 0.01 {
		fmt.Printf("calculating f_n(z,c) = z^%.2f + %.2f*z^2 + c\n", n, n)
		img := image.NewPaletted(image.Rect(0, 0, width, height), palette)

		images = append(images, img)
		delays = append(delays, 0)

		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				// map pixel location to point on complex plane.
				pt := mandelbrot(
					complex(
						float64(x)/scale+rMin,
						float64(y)/scale+iMin,
					),
					n,
				)

				img.Set(x, y, getColor(pt))
			}
		}
	}

	f, err := os.OpenFile("mandelbrot-complex-polynomial-2.gif", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	if err = gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	}); err != nil {
		fmt.Println(err)
		return
	}
}
