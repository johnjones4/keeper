package api

import (
	"main/core"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type API struct {
	runtime *core.RuntimeContext
	r       *chi.Mux
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.r.ServeHTTP(w, r)
}

func New(runtime *core.RuntimeContext) *API {
	a := &API{
		runtime: runtime,
		r:       chi.NewRouter(),
	}

	a.r.Use(middleware.RequestID)
	a.r.Use(middleware.RealIP)
	a.r.Use(middleware.Logger)
	a.r.Use(middleware.Recoverer)

	a.r.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		r.Post("/token", a.newToken)

		r.Route("/note", func(r chi.Router) {
			r.Use(a.verifyToken)

			r.Get("/", a.getNotes)
			r.Post("/", a.newNote)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", a.getNote)
				r.Put("/", a.updateNote)
			})
		})
	})

	return a
}
