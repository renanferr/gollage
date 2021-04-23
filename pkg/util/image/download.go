package image

import (
	"errors"
	"image"
	"log"
	"net/http"

	"github.com/renanferr/gollage/pkg/albums"
	"golang.org/x/sync/errgroup"
)

type Downloader struct {
	albums  []*albums.Album
	decoder *Decoder
}

func NewDownloader(albums []*albums.Album) *Downloader {
	d := new(Downloader)
	d.albums = albums
	d.decoder = NewDecoder()
	return d
}

func (d *Downloader) DownloadAll() {
	var g errgroup.Group
	images := make(map[int]image.Image)
	for _, a := range d.albums {
		url := a.ImageUrl
		rank := a.Rank
		g.Go(func() error {
			img, err := d.download(url)
			if err != nil {
				log.Printf("error downloading image: %s", err.Error())
			}
			images[rank] = img
			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		log.Printf("\nerror in waiting group: %s", err)
	}

	for _, a := range d.albums {
		a.CoverImage = images[a.Rank]
	}

}

func (d *Downloader) download(url string) (image.Image, error) {
	log.Printf("downloading image: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()
	img, err := d.decoder.Decode(resp.Body, d.decoder.Format(url))
	if errors.Is(err, image.ErrFormat) {
		log.Printf("error decoding image: %s", err.Error())
		return nil, nil
	}
	log.Println("downloaded image successfully")
	return img, err
}
