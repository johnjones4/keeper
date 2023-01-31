package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
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

func errorResponse(w http.ResponseWriter, status int, err error) {
	log.Print(err)
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
