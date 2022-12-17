package api

import (
	"main/core"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *API) putNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	dbNote, err := a.runtime.Store.GetNote(id)
	if err != nil {
		errorResponse(a.runtime.Log, w, 0, err)
		return
	}

	var inputNote core.Note
	err = readJson(r, &inputNote)
	if err != nil {
		errorResponse(a.runtime.Log, w, 0, err)
		return
	}

	oldName := dbNote.Title
	dbNote.Body = inputNote.Body
	dbNote.Title = inputNote.Title

	for _, p := range a.runtime.Processors {
		err = p(&dbNote)
		if err != nil {
			errorResponse(a.runtime.Log, w, 0, err)
			return
		}
	}

	if oldName == dbNote.Title {
		err = a.runtime.Store.SaveNote(&dbNote)
	} else {
		newName := dbNote.Title
		dbNote.Title = oldName
		err = a.runtime.Store.RenameAndSaveNote(&dbNote, newName)
	}
	if err != nil {
		errorResponse(a.runtime.Log, w, 0, err)
		return
	}

	jsonResponse(w, http.StatusOK, dbNote)
}
