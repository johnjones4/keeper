package main

import (
	"main/api"
	"main/core"
	"main/store"
	"net/http"
	"os"
)

func main() {
	rc := core.RuntimeContext{
		Store:        store.New(os.Getenv("ROOT_DIR")),
		PrivateKey:   []byte(os.Getenv("PRIVATE_KEY")),
		PasswordHash: []byte(os.Getenv("PASSWORD")),
	}

	a := api.New(&rc)

	err := http.ListenAndServe(os.Getenv("HTTP_HOST"), a)
	panic(err)
}
