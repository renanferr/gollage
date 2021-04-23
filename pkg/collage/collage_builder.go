package collage

import (
	"context"

	"github.com/renanferr/gollage/pkg/albums"
	"golang.org/x/sync/errgroup"
)

type AsyncDownloader interface {
	AsyncDownloadAll(ctx context.Context, out chan *albums.Album) error
}

type CollageBuilder struct {
	albums     []*albums.Album
	downloader AsyncDownloader
}

func NewCollageBuilder() *CollageBuilder {
	return &CollageBuilder{}
}

func (cb *CollageBuilder) WithDownloader(d AsyncDownloader) *CollageBuilder {
	cb.downloader = d
	return cb
}

func (cb *CollageBuilder) WithAlbums(a []*albums.Album) *CollageBuilder {
	cb.albums = a
	return cb
}

func (cb *CollageBuilder) Build(ctx context.Context) (*Collage, error) {
	c := new(Collage)
	ch := make(chan *albums.Album, len(cb.albums))
	errG, ctx := errgroup.WithContext(ctx)

	errG.Go(func() error {
		return cb.download(ctx, ch)
	})

	errG.Go(func() error {
		return cb.compose(ctx, ch)
	})

	if err := errG.Wait(); err != nil {
		return nil, err
	}

	return c, nil
}

func (cb *CollageBuilder) download(ctx context.Context, out chan *albums.Album) error {
	return cb.downloader.AsyncDownloadAll(ctx, out)
}

func (cb *CollageBuilder) compose(ctx context.Context, in chan *albums.Album) error {
	return nil
}
