package albums

import (
	"image"
	"log"
	"strconv"

	"github.com/shkh/lastfm-go/lastfm"
)

type Album struct {
	Rank       int         `json:"rank"`
	Name       string      `json:"name"`
	Artist     string      `json:"artist"`
	ImageUrl   string      `json:"imageUrl"`
	CoverImage image.Image `json:"coverImage"`
}

func AlbumsFromLastFM(result lastfm.UserGetTopAlbums) ([]*Album, error) {
	albums := []*Album{}
	images, err := ImagesFromLastFM(result)

	if err != nil {
		return nil, err
	}

	for _, album := range result.Albums {
		a := new(Album)

		rank, err := strconv.Atoi(album.Rank)
		if err != nil {
			return nil, err
		}
		a.Rank = rank

		a.Name = album.Name
		a.Artist = album.Artist.Name

		if err != nil {
			return nil, err
		}
		a.ImageUrl = images[album.Rank]
		log.Printf("\n%s) %s - %s", album.Rank, album.Name, album.Artist.Name)
		albums = append(albums, a)
	}
	return albums, nil
}

func ImagesFromLastFM(result lastfm.UserGetTopAlbums) (map[string]string, error) {
	images := make(map[string]string)

	for _, a := range result.Albums {
		for _, i := range a.Images {
			if i.Size == "extralarge" {
				images[a.Rank] = i.Url
			}
		}
	}
	return images, nil
}
