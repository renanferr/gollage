package collage

import (
	"image"
	"image/color"
)

type Collage struct {
	*image.RGBA
}

func NewCollage(img *image.RGBA) *Collage {
	return &Collage{img}
}

func (c *Collage) FillAll(color color.Color) {
	c.Fill(color, c.Bounds().Min.X, c.Bounds().Min.Y, c.Bounds().Max.X, c.Bounds().Max.Y)
}

func (c *Collage) Fill(color color.Color, x0, y0, x1, y1 int) {
	for x := x0; x <= x1; x++ {
		for y := y0; y <= y1; y++ {
			c.Set(x, y, color)
		}
	}
}

func (c *Collage) Paste(img image.Image, x, y int) {
	x0 := x
	// y0 := y
	for w := img.Bounds().Min.Y; w <= img.Bounds().Max.Y; w++ {
		for z := img.Bounds().Min.X; z <= img.Bounds().Max.X; z++ {
			c.Set(x, y, img.At(z, w))
			x++
		}
		x = x0
		y++
	}
}
