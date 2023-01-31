package store

import (
	"errors"
	"main/core"
	"os"
	"path"
)

func (s *StaticFileStore) SaveNote(n *core.Note, failOnOverwrite bool, failOnNew bool) error {
	n.Key = path.Clean(n.Key)

	notePath := path.Join(s.rootDir, n.Key)

	stat, err := os.Stat(notePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if stat != nil && failOnOverwrite {
		return core.ErrorAlreadyExists
	}
	if stat == nil && failOnNew {
		return core.ErrorDoesNotExist
	}

	dir := path.Dir(notePath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(notePath, []byte(n.Body), 0755)
	if err != nil {
		return err
	}

	stat, err = os.Stat(notePath)
	if err != nil {
		return err
	}

	n.Modified = stat.ModTime().UTC()

	return nil
}
