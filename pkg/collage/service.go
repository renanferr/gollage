package collage

import (
	"image"
	"image/color"
	"log"

	"github.com/renanferr/gollage/pkg/albums"
)

var RGBA_BLACK = color.RGBA{0x00, 0x00, 0x00, 0x00}

type Service interface {
	Compose(albums []*albums.Album, rows, cols, width, height int) image.Image
}

type service struct{}

func New() Service {
	return &service{}
}

func getBounds(img image.Image) (int, int, int, int) {
	return img.Bounds().Min.X,
		img.Bounds().Min.Y,
		img.Bounds().Max.X,
		img.Bounds().Max.Y
}

func (s *service) Compose(albums []*albums.Album, rows, cols, width, height int) image.Image {
	rect := image.Rect(0, 0, width, height)
	rgba := image.NewRGBA(rect)
	collage := NewCollage(rgba)

	offsetX, offsetY := collage.Rect.Size().X/rows,
		collage.Rect.Size().Y/cols

	i := 0
	x, y := 0, 0

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if i >= len(albums) || albums[i].CoverImage == nil {
				collage.Fill(color.RGBA{}, x, y, x+offsetX, y+offsetY)
				log.Printf("skipped image #%d", i)
				x += offsetX
				i++
				continue
			}
			collage.Paste(albums[i].CoverImage, x, y)
			log.Printf("pasted image #%d", i)
			x += offsetX
			i++
		}
		y += offsetY
		x = 0
	}
	// return images[2]
	return collage
}
