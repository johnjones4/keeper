package api

import (
	"main/core"
	"net/http"
)

func (a *API) newNote(w http.ResponseWriter, r *http.Request) {
	var note core.Note

	err := readJson(r, &note)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	err = a.runtime.Store.SaveNote(&note, true, false)
	if err != nil {
		if err == core.ErrorAlreadyExists {
			errorResponse(w, http.StatusBadRequest, err)
		} else {
			errorResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

	jsonResponse(w, http.StatusOK, note)
}
