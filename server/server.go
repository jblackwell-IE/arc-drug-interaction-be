package server

import (
	"net/http"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/IdeaEvolver/cutter-pkg/service"
	"github.com/go-chi/chi"
	"github.com/jblackwell-IE/arc-drug-interaction-be/fdb"
	"github.com/rs/cors"
	"go.opencensus.io/plugin/ochttp"
)

type Handler struct {
	Interactions *fdb.Client
}

func New(cfg *service.Config, handler *Handler) *service.Server {
	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodHead, http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPost},
		ExposedHeaders:   []string{"Authorization"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler)

	router.Handle("/*", http.FileServer(http.Dir("./dist")))
	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		return
	})

	router.Route("/api/v1", func(router chi.Router) {
		router.Route("/", func(router chi.Router) {
			router.Method("POST", "/check-interactions", service.JsonHandler(handler.GetDrugInteractions))
		})
	})

	httpHandler := &ochttp.Handler{
		// Use the Google Cloud propagation format.
		Propagation: &propagation.HTTPFormat{},
		Handler:     router,
	}

	return service.GracefulServer(cfg, httpHandler)
}
