package store

import (
	"main/core"
	"os"
	"path"
)

func (s *StaticFileStore) GetNote(key string) (core.Note, error) {
	cKey := path.Clean(key)

	notePath := path.Join(s.rootDir, cKey)
	if !s.isPathSafe(notePath) {
		return core.Note{}, core.ErrorBadPath
	}

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

	n.Modified = core.NoteTime(stat.ModTime().UTC())

	return n, err
}
