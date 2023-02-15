package core

import "errors"

var (
	ErrorAlreadyExists = errors.New("note already exists")
	ErrorDoesNotExist  = errors.New("note does not exist")
	ErrorBadPath       = errors.New("not a valid path")
)
