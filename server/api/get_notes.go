package api

import (
	"net/http"
	"net/url"
)

func (a *API) GetNote(w http.ResponseWriter, r *http.Request, params GetNoteParams) {
	if !a.verifyToken(w, r) {
		return
	}

	var notes []string
	var nextPage string
	var err error

	if params.Q != nil {
		dir, err := url.QueryUnescape(*params.Q)
		if err != nil {
			a.errorResponse(w, http.StatusBadRequest, err)
			return
		}
		notes, err = a.runtime.Index.Search(dir)
		if err != nil {
			a.errorResponse(w, http.StatusInternalServerError, err)
			return
		}
	} else if params.Dir != nil {
		notes, err = a.runtime.Store.GetNoteDirectory(*params.Dir)
		if err != nil {
			a.errorResponse(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		var page string
		if params.Page != nil {
			page = *params.Page
		}
		notes, nextPage, err = a.runtime.Store.GetNotes(3, page)
		if err != nil {
			a.errorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	resp := Notes{
		Notes:    notes,
		NextPage: &nextPage,
	}

	a.jsonResponse(w, http.StatusOK, resp)
}
