package api

import (
	"errors"
	"main/core"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

func (a *API) updateNote(w http.ResponseWriter, r *http.Request) {
	var note core.Note

	err := readJson(r, &note)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	id, err := url.QueryUnescape(chi.URLParam(r, "id"))
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	key, err := unmarshalKey(id)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	if key != note.Key {
		errorResponse(w, http.StatusBadRequest, errors.New("key mismatch"))
		return
	}

	err = a.runtime.Store.SaveNote(&note, false, true)
	if err != nil {
		if err == core.ErrorDoesNotExist {
			errorResponse(w, http.StatusBadRequest, err)
		} else {
			errorResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

	jsonResponse(w, http.StatusOK, note)
}
