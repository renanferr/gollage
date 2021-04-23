package albums

import "github.com/shkh/lastfm-go/lastfm"

type Service interface {
	GetTopAlbums(user string, period string, limit int) ([]*Album, error)
}

type service struct {
	api *lastfm.Api
}

func New(api *lastfm.Api) *service {
	return &service{
		api: api,
	}
}

func (s *service) GetTopAlbums(user string, period string, limit int) ([]*Album, error) {
	resp, err := s.api.User.GetTopAlbums(lastfm.P{"user": user, "period": period, "limit": limit})
	if err != nil {
		return nil, err
	}

	return AlbumsFromLastFM(resp)
}
