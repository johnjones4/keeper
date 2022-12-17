package hybridstore

import (
	"main/core"
)

func (s *HybridStore) GetNote(id string) (core.Note, error) {
	row := s.db.QueryRow("SELECT id, path, title, sourceURL, source, format, created, updated FROM notes WHERE id = $1", id)

	note, err := parseRow(row)
	if err != nil {
		return core.Note{}, err
	}

	err = s.populateTags(&note)
	if err != nil {
		return core.Note{}, err
	}

	err = s.loadNoteBody(&note)
	if err != nil {
		return core.Note{}, err
	}

	return note, nil
}
