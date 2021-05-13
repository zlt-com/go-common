package common

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"os"
)

// Circle 切成圆形
func Circle(sourceImg, newImg string) {
	file, err := os.Create(newImg)

	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	imageFile, err := os.Open(sourceImg)

	if err != nil {
		fmt.Println(err)
	}
	defer imageFile.Close()

	srcImg, _ := jpeg.Decode(imageFile)

	w := srcImg.Bounds().Max.X - srcImg.Bounds().Min.X
	h := srcImg.Bounds().Max.Y - srcImg.Bounds().Min.Y

	d := w
	if w > h {
		d = h
	}

	dstImg := NewCircleMask(srcImg, image.Point{d / (w), d / (w)}, d-1)

	png.Encode(file, dstImg)
}

// NewCircleMask NewCircleMask
func NewCircleMask(img image.Image, p image.Point, d int) CircleMask {
	return CircleMask{img, p, d}
}

// CircleMask CircleMask
type CircleMask struct {
	image    image.Image
	point    image.Point
	diameter int
}

// ColorModel ColorModel
func (ci CircleMask) ColorModel() color.Model {
	return ci.image.ColorModel()
}

// Bounds Bounds
func (ci CircleMask) Bounds() image.Rectangle {
	return image.Rect(0, 0, ci.diameter, ci.diameter)
}

// At At
func (ci CircleMask) At(x, y int) color.Color {
	d := ci.diameter
	dis := math.Sqrt(math.Pow(float64(x-d/2), 2) + math.Pow(float64(y-d/2), 2))
	if dis > float64(d)/2 {
		return ci.image.ColorModel().Convert(color.RGBA{0, 0, 0, 1})
	}
	return ci.image.At(ci.point.X+x, ci.point.Y+y)
}
