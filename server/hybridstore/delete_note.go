package hybridstore

import (
	"main/core"
	"os"
	"path"
)

func (s *HybridStore) DeleteNote(note *core.Note) error {
	_, err := s.db.Exec("DELETE FROM notes WHERE id = $1", note.ID)
	if err != nil {
		return err
	}

	oldFullPath := path.Join(s.filePath, note.Path)
	newFullPath := path.Join(s.recyclePath, note.Path)

	err = os.Rename(oldFullPath, newFullPath)
	if err != nil {
		return err
	}

	return nil
}
