package hybridstore

import (
	"os"
	"path"

	"github.com/johnjones4/keeper/core"
)

func (s *HybridStore) DeleteNote(note *core.Note) error {
	_, err := s.db.Exec("DELETE FROM notes WHERE id = $1", note.ID)
	if err != nil {
		return err
	}

	oldFullPath := path.Join(s.filePath, note.Path)
	newFullPath := path.Join(s.recyclePath, note.Path)

	newFullPathDir := path.Dir(newFullPath)
	err = os.MkdirAll(newFullPathDir, 0755)
	if err != nil {
		return err
	}

	err = os.Rename(oldFullPath, newFullPath)
	if err != nil {
		return err
	}

	return nil
}
