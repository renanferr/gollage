package collage

import (
	"context"

	"github.com/renanferr/gollage/pkg/albums"
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

func (cb *CollageBuilder) Build(ctx context.Context) *Collage {
	c := new(Collage)
	// out := make(chan *albums.Album, len(cb.albums))
	// errG, ctx := errgroup.WithContext(ctx)
	return c
}
