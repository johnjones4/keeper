package hybridstore

import "github.com/johnjones4/keeper/core"

func (s *HybridStore) GetNoteByPath(path string) (core.Note, error) {
	row := s.db.QueryRow("SELECT id, path, title, sourceURL, source, format, created, updated FROM notes WHERE path = $1", path)

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
