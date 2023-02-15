package store

func (s *StaticFileStore) GetNoteDirectory(directory string) ([]string, error) {
	col := collector{
		notes:             make([]string, 0),
		pageSize:          -1,
		hasFoundStartPath: true,
	}

	_, err := s.exploreDir(&col, "", directory, 0, 1, true)
	if err != nil {
		return nil, err
	}

	return col.notes, nil
}
