package main

import (
	"main/api"
	"main/core"
	"main/hybridstore"
	"main/processors"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	hstore, err := hybridstore.New(
		os.Getenv("DB_FILE"),
		os.Getenv("DOCS_DIR"),
		os.Getenv("TRASH_DIR"),
	)
	if err != nil {
		log.Panic(err)
	}

	runtime := core.Runtime{
		Store: hstore,
		Log:   log,
	}

	processorsI := processors.Processors{
		Runtime: &runtime,
	}
	runtime.Processors = processorsI.All()

	h := api.New(&runtime)
	err = http.ListenAndServe(os.Getenv("HTTP_HOST"), h)
	if err != nil {
		log.Panic(err)
	}
}
