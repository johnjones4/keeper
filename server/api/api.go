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

	a.r.Group(func(mux chi.Router) {
		options := ChiServerOptions{
			BaseURL:    "/api",
			BaseRouter: mux,
			ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				runtime.Log.Error(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
			},
		}
		HandlerWithOptions(a, options)
	})

	return a
}
