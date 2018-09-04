package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func importRGB(path string) image.RGBA {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()
	dest := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.RGBAModel.Convert(img.At(x, y))
			rgb, _ := c.(color.RGBA)
			dest.Set(x, y, rgb)
		}
	}

	return *dest
}

func filtering(img image.RGBA) image.RGBA {
	filter := [][]float64{
		{1, 1, 1},
		{1, 1, 1},
		{1, 1, 1}}

	// * 1/16
	bounds := img.Bounds()
	filtered := *image.NewRGBA(bounds)

	for y := bounds.Min.Y + 1; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X; x++ {
			var r, g, b float64

			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					r += float64(img.RGBAAt(x+dx, y+dy).R) * filter[dy+1][dx+1]
					g += float64(img.RGBAAt(x+dx, y+dy).G) * filter[dy+1][dx+1]
					b += float64(img.RGBAAt(x+dx, y+dy).B) * filter[dy+1][dx+1]
				}
			}
			newPix := color.RGBA{R: uint8(r / 9), G: uint8(g / 9), B: uint8(b / 9), A: 255}
			filtered.SetRGBA(x, y, newPix)
		}
	}

	return filtered
}

func main() {
	img := importRGB("../data/lenna.png")

	filterd := filtering(img)

	file, err := os.Create("../data/gaussian-filter.png")
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(file, &filterd)
	fmt.Println("export gaussian-filterd image")
}
