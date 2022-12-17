package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *API) getNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	note, err := a.runtime.Store.GetNote(id)
	if err != nil {
		errorResponse(a.runtime.Log, w, 0, err)
		return
	}

	jsonResponse(w, http.StatusOK, note)
}
