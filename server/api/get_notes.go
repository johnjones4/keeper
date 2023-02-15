package api

import (
	"net/http"
)

type notesResponse struct {
	Notes    []string `json:"notes"`
	NextPage string   `json:"nextPage"`
}

func (a *API) getNotes(w http.ResponseWriter, r *http.Request) {
	var notes []string
	var nextPage string
	var err error

	if r.URL.Query().Has("q") {
		q := r.URL.Query().Get("q")
		notes, err = a.runtime.Index.Search(q)
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, err)
			return
		}
	} else if r.URL.Query().Has("dir") {
		notes, err = a.runtime.Store.GetNoteDirectory(r.URL.Query().Get("dir"))
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		page := r.URL.Query().Get("page")
		notes, nextPage, err = a.runtime.Store.GetNotes(3, page)
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	resp := notesResponse{
		Notes:    notes,
		NextPage: nextPage,
	}

	jsonResponse(w, http.StatusOK, resp)
}
