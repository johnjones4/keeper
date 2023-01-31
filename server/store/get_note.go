package store

import (
	"main/core"
	"os"
	"path"
)

func (s *StaticFileStore) GetNote(key string) (core.Note, error) {
	cKey := path.Clean(key)

	notePath := path.Join(s.rootDir, cKey)

	contents, err := os.ReadFile(notePath)
	if err != nil {
		return core.Note{}, err
	}

	n := core.Note{
		Key:  cKey,
		Body: string(contents),
	}

	stat, err := os.Stat(notePath)
	if err != nil {
		return core.Note{}, err
	}

	n.Modified = stat.ModTime().UTC()

	return n, err
}
