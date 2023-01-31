package api

import "net/http"

type notesResponse struct {
	Notes    []string `json:"notes"`
	NextPage string   `json:"nextPage"`
}

func (a *API) getNotes(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")

	notes, nextPage, err := a.runtime.Store.GetNotes(3, page)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	resp := notesResponse{
		Notes:    notes,
		NextPage: nextPage,
	}

	jsonResponse(w, http.StatusOK, resp)
}
