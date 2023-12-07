package api

import (
	"main/core"
	"net/http"
)

func (a *API) newNote(w http.ResponseWriter, r *http.Request) {
	var note core.Note

	err := readJson(r, &note)
	if err != nil {
		a.errorResponse(w, http.StatusBadRequest, err)
		return
	}

	err = a.runtime.Store.SaveNote(&note, true, false)
	if err != nil {
		if err == core.ErrorAlreadyExists {
			a.errorResponse(w, http.StatusBadRequest, err)
		} else {
			a.errorResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

	err = a.runtime.Index.Add(&note)
	if err != nil {
		a.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	a.jsonResponse(w, http.StatusOK, note)
}
