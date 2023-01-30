package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/johnjones4/keeper/core"
)

func initNote(title string) (core.Note, error) {
	urlRoot := os.Getenv("KEEPER_URL_ROOT")
	if urlRoot == "" {
		return core.Note{}, errors.New("env variable KEEPER_URL_ROOT not set")
	}
	noteIn := core.Note{
		Title:  title,
		Format: "text/markdown",
	}
	bodyIn, err := json.Marshal(noteIn)
	if err != nil {
		return core.Note{}, err
	}

	res, err := http.Post(fmt.Sprintf("%s/api/note", urlRoot), "application/json", bytes.NewBuffer(bodyIn))
	if err != nil {
		return core.Note{}, err
	}

	bodyOut, err := io.ReadAll(res.Body)
	if err != nil {
		return core.Note{}, err
	}

	var noteOut core.Note
	err = json.Unmarshal(bodyOut, &noteOut)
	if err != nil {
		return core.Note{}, err
	}

	return noteOut, nil
}

func finalPathFromNote(note core.Note) (string, error) {
	basePath := os.Getenv("KEEPER_PATH")
	if basePath == "" {
		return "", errors.New("env variable KEEPER_PATH not set")
	}
	return path.Join(basePath, note.Path), nil
}

func main() {
	if len(os.Args) != 2 {
		panic("Usage: note [title]")
	}
	title := os.Args[1]

	note, err := initNote(title)
	if err != nil {
		panic(err)
	}

	notePath, err := finalPathFromNote(note)
	if err != nil {
		panic(err)
	}

	_, err = os.Create(notePath)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("/usr/bin/vim", notePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
