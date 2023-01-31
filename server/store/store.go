package store

type StaticFileStore struct {
	rootDir string
}

func New(rootDir string) *StaticFileStore {
	return &StaticFileStore{rootDir}
}
