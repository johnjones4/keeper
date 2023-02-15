package store

import "path"

type StaticFileStore struct {
	rootDir string
	ignore  []string
}

func New(rootDir string, ignore []string) *StaticFileStore {
	return &StaticFileStore{path.Clean(rootDir), ignore}
}
