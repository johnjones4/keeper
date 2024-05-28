package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"main/core"
	"net/http"
	"time"
)

const (
	issuer     = "keeper"
	expiration = time.Hour * 24 * 30
)

func unmarshalKey(id string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (a *API) jsonResponse(w http.ResponseWriter, status int, info any) {
	bytes, err := json.Marshal(info)
	if err != nil {
		a.runtime.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	w.Write(bytes)
}

func (a *API) errorResponse(w http.ResponseWriter, status int, err error) {
	a.runtime.Log.Error(err)
	msg := map[string]any{
		"ok":      false,
		"message": err.Error(),
	}
	a.jsonResponse(w, status, msg)
}

func readJson(r *http.Request, readTo any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, readTo)
	if err != nil {
		return err
	}

	return nil
}

func mapFromCoreNote(note core.Note) Note {
	return Note{
		Key:      note.Key,
		Body:     note.Body,
		Modified: (*time.Time)(&note.Modified),
	}
}

func mapToCoreNote(note Note) core.Note {
	now := time.Now()
	if note.Modified == nil {
		note.Modified = &now
	}
	return core.Note{
		Key:      note.Key,
		Body:     note.Body,
		Modified: core.NoteTime(*note.Modified),
	}
}
