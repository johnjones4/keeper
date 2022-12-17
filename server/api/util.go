package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func jsonResponse(w http.ResponseWriter, status int, info any) {
	bytes, err := json.Marshal(info)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	w.Write(bytes)
}

func errorResponse(log logrus.FieldLogger, w http.ResponseWriter, status int, err error) {
	log.Error(err)

	if status == 0 {
		status = mapErrorToStatus(err)
	}

	msg := map[string]any{
		"ok":      false,
		"message": err.Error(),
	}
	jsonResponse(w, status, msg)
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

func mapErrorToStatus(err error) int {
	switch err {
	case sql.ErrNoRows:
		return http.StatusNotFound
	}
	switch err.(type) {
	case *json.InvalidUnmarshalError:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func parseDateQuery(ds string) time.Time {
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		time.Stamp,
		"2006-01-02 15:04:05",
	}
	for _, format := range formats {
		t, err := time.Parse(format, ds)
		if err == nil {
			return t
		}
	}
	return time.Time{}
}
