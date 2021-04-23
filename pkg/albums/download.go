package albums

import (
	"context"
	"errors"
	"image"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type Downloader struct {
	albums  []*Album
	decoder *Decoder
}

func NewDownloader() *Downloader {
	d := new(Downloader)
	d.albums = []*Album{}
	d.decoder = NewDecoder()
	return d
}

func (d *Downloader) WithAlbums(a ...*Album) *Downloader {
	d.albums = append(d.albums, a...)
	return d
}

func (d *Downloader) AsyncDownloadAll(ctx context.Context, out chan *Album) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, a := range d.albums {
		func(a *Album) {
			g.Go(func() error {
				img, err := d.download(ctx, a.ImageUrl)
				if err != nil {
					log.Printf("error downloading image: %s", err.Error())
				}
				a.CoverImage = img
				out <- a
				return nil
			})

		}(a)
	}
	return g.Wait()

}

func (d *Downloader) DownloadAll(ctx context.Context) {
	var g errgroup.Group
	images := make(map[int]image.Image)
	for _, a := range d.albums {
		url := a.ImageUrl
		rank := a.Rank
		g.Go(func() error {
			img, err := d.download(ctx, url)
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

func (d *Downloader) download(ctx context.Context, url string) (image.Image, error) {
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
