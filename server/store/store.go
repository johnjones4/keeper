package store

type StaticFileStore struct {
	rootDir string
	ignore  []string
}

func New(rootDir string, ignore []string) *StaticFileStore {
	return &StaticFileStore{rootDir, ignore}
}
