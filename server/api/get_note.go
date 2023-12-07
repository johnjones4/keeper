package api

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

func (a *API) getNote(w http.ResponseWriter, r *http.Request) {
	id, err := url.QueryUnescape(chi.URLParam(r, "id"))
	if err != nil {
		a.errorResponse(w, http.StatusBadRequest, err)
		return
	}

	key, err := unmarshalKey(id)
	if err != nil {
		a.errorResponse(w, http.StatusBadRequest, err)
		return
	}

	note, err := a.runtime.Store.GetNote(key)
	if err != nil {
		a.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	a.jsonResponse(w, http.StatusOK, note)
}
