package api

import (
	"errors"
	"main/core"
	"net/http"
	"net/url"
)

func (a *API) PutNoteKey(w http.ResponseWriter, r *http.Request, dirtyId string, params PutNoteKeyParams) {
	if !a.verifyToken(w, r) {
		return
	}

	var inNote Note

	err := readJson(r, &inNote)
	if err != nil {
		a.errorResponse(w, http.StatusBadRequest, err)
		return
	}

	id, err := url.QueryUnescape(dirtyId)
	if err != nil {
		a.errorResponse(w, http.StatusBadRequest, err)
		return
	}

	key, err := unmarshalKey(id)
	if err != nil {
		a.errorResponse(w, http.StatusBadRequest, err)
		return
	}

	note := mapToCoreNote(inNote)

	if key != note.Key {
		a.errorResponse(w, http.StatusBadRequest, errors.New("key mismatch"))
		return
	}

	err = a.runtime.Store.SaveNote(&note, false, true)
	if err != nil {
		if err == core.ErrorDoesNotExist {
			a.errorResponse(w, http.StatusBadRequest, err)
		} else {
			a.errorResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

	err = a.runtime.Index.Update(&note)
	if err != nil {
		a.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	a.jsonResponse(w, http.StatusOK, mapFromCoreNote(note))
}
