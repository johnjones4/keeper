package api

import (
	"main/core"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type API struct {
	runtime *core.Runtime
	handler *chi.Mux
}

func New(runtime *core.Runtime) *API {
	h := API{
		runtime: runtime,
		handler: chi.NewRouter(),
	}

	h.handler.Use(middleware.RequestID)
	h.handler.Use(middleware.RealIP)
	h.handler.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: runtime.Log, NoColor: false}))
	h.handler.Use(middleware.Recoverer)

	h.handler.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		r.Route("/note", func(r chi.Router) {
			r.Get("/", h.getNotes)
			r.Post("/", h.postNote)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.getNote)
				r.Put("/", h.putNote)
			})
		})
	})

	return &h
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(w, r)
}
