package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *API) deleteNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	dbNote, err := a.runtime.Store.GetNote(id)
	if err != nil {
		errorResponse(a.runtime.Log, w, 0, err)
		return
	}

	err = a.runtime.Store.DeleteNote(&dbNote)
	if err != nil {
		errorResponse(a.runtime.Log, w, 0, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
