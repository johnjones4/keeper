package hybridstore

import (
	"main/core"
	"os"
	"path"
)

func (s *HybridStore) RenameAndSaveNote(note *core.Note, newName string) error {
	oldPath := note.Path
	noteDir := path.Dir(oldPath)

	note.Title = newName
	filename, err := generateNoteFilename(newName, note.Format)
	if err != nil {
		return err
	}
	note.Path = path.Join(noteDir, filename) //TODO format change

	fullOldPath := path.Join(s.filePath, oldPath)
	fullNewPath := path.Join(s.filePath, note.Path)

	err = os.Rename(fullOldPath, fullNewPath)
	if err != nil {
		return err
	}

	return s.SaveNote(note)
}
