package main

import (
	"main/api"
	"main/core"
	"main/index"
	"main/store"
	"net/http"
	"os"
	"strings"
	"time"
)

func startIndex(idx core.Index) {
	for {
		err := idx.ReIndex()
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Hour)
	}
}

func main() {
	store := store.New(os.Getenv("ROOT_DIR"), strings.Split(os.Getenv("IGNORE_FS"), "|"))
	idx, err := index.New(os.Getenv("INDEX_PATH"), store)
	if err != nil {
		panic(err)
	}

	rc := core.RuntimeContext{
		Store:        store,
		Index:        idx,
		PrivateKey:   []byte(os.Getenv("PRIVATE_KEY")),
		PasswordHash: []byte(os.Getenv("PASSWORD")),
	}

	go startIndex(idx)

	a := api.New(&rc)
	err = http.ListenAndServe(os.Getenv("HTTP_HOST"), a)
	panic(err)
}
