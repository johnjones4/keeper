package api

import (
	"main/core"
	"net/http"
)

func (a *API) PostNote(w http.ResponseWriter, r *http.Request, params PostNoteParams) {
	if !a.verifyToken(w, r) {
		return
	}

	var inNote Note

	err := readJson(r, &inNote)
	if err != nil {
		a.errorResponse(w, http.StatusBadRequest, err)
		return
	}

	note := mapToCoreNote(inNote)

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

	a.jsonResponse(w, http.StatusOK, mapFromCoreNote(note))
}
