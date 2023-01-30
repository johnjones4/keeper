package api

import (
	"net/http"

	"github.com/johnjones4/keeper/core"
)

func (a *API) postNote(w http.ResponseWriter, r *http.Request) {
	var note core.Note
	err := readJson(r, &note)
	if err != nil {
		errorResponse(a.runtime.Log, w, 0, err)
		return
	}

	for _, p := range a.runtime.Processors {
		err = p(&note)
		if err != nil {
			errorResponse(a.runtime.Log, w, 0, err)
			return
		}
	}

	err = a.runtime.Store.SaveNote(&note)
	if err != nil {
		errorResponse(a.runtime.Log, w, 0, err)
		return
	}

	jsonResponse(w, http.StatusOK, note)
}
