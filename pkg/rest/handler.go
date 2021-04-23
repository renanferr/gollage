package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/renanferr/gollage/pkg/albums"
	"github.com/renanferr/gollage/pkg/collage"
	collageApi "github.com/renanferr/gollage/pkg/rest/collage"
)

const ApiPrefix = "/api"

func Handler(c collage.Service, a albums.Service) http.Handler {
	apiRouter := chi.NewRouter()

	apiRouter.Use(middleware.Recoverer, middleware.RequestID, middleware.Logger)
	apiRouter.Mount("/collage", collageApi.Handler(c, a))

	router := chi.NewRouter()

	router.Get("/healthcheck", handleHealthcheck)

	router.Mount(ApiPrefix, apiRouter)

	return router
}

func handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"alive": true})
}
