package api

import (
	"net/http"

	"github.com/johnjones4/keeper/core"
)

func (a *API) getNotes(w http.ResponseWriter, r *http.Request) {
	var query core.NotesQuery

	if r.URL.Query().Has("start") {
		query.Start = parseDateQuery(r.URL.Query().Get("start"))
	}

	if r.URL.Query().Has("end") {
		query.End = parseDateQuery(r.URL.Query().Get("end"))
	}

	if r.URL.Query().Has("text") {
		query.Text = r.URL.Query().Get("text")
	}

	if r.URL.Query().Has("tag") {
		query.Tags = r.URL.Query()["tag"]
	}

	notes, err := a.runtime.Store.GetNotes(query)
	if err != nil {
		errorResponse(a.runtime.Log, w, 0, err)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]any{
		"items": notes,
	})
}
