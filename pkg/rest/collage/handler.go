package collage

import (
	"encoding/json"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/renanferr/gollage/pkg/albums"
	"github.com/renanferr/gollage/pkg/collage"
)

const (
	username   = "ohships"
	period     = "7day"
	rows, cols = 2, 2
	limit      = rows * cols
	width      = rows * 300
	height     = cols * 300
)

type Options struct {
	Rows, Cols, Width, Height, Limit int
	Username, Period                 string
}

func OptionsFromRequest(r *http.Request) (*Options, error) {
	q := r.URL.Query()
	rows, err := strconv.Atoi(getWithDefault(q, "rows", "3"))
	if err != nil {
		return nil, err
	}
	cols, err := strconv.Atoi(getWithDefault(q, "cols", "3"))
	if err != nil {
		return nil, err
	}

	return &Options{
		Rows:     rows,
		Cols:     cols,
		Width:    rows * 300,
		Height:   cols * 300,
		Limit:    rows * cols,
		Username: q.Get("username"),
		Period:   getWithDefault(q, "period", "7day"),
	}, nil
}

func getWithDefault(q url.Values, param string, defaultValue string) string {
	v := q.Get(param)
	if v == "" {
		return defaultValue
	}
	return v
}

func Handler(c collage.Service, a albums.Service) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", getCollage(c, a))

	return r
}

type errorResponse struct {
	Message string `json:"message"`
}

func sendErrorResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&errorResponse{message}); err != nil {
		log.Panicf("error encoding error response: %s", err)
	}
}

func getCollage(c collage.Service, a albums.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		opts, err := OptionsFromRequest(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		a, err := a.GetTopAlbums(opts.Username, opts.Period, opts.Limit)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// d.DownloadAll()

		// log.Println(opts)
		img := c.Compose(r.Context(), a, opts.Rows, opts.Cols, opts.Width, opts.Height)
		w.Header().Set("Content-type", "image/png")
		png.Encode(w, img)
		w.WriteHeader(http.StatusOK)
	}
}
