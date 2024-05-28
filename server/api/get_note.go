package api

import (
	"net/http"
	"net/url"
)

func (a *API) GetNoteKey(w http.ResponseWriter, r *http.Request, dirtyId string, params GetNoteKeyParams) {
	if !a.verifyToken(w, r) {
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

	note, err := a.runtime.Store.GetNote(key)
	if err != nil {
		a.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	a.jsonResponse(w, http.StatusOK, mapFromCoreNote(note))
}
