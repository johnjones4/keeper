package store

import (
	"encoding/base64"
)

func (s *StaticFileStore) GetNotes(pageSize int, page string) ([]string, string, error) {
	startPath := ""
	started := true
	if page != "" {
		startPathBytes, err := base64.StdEncoding.DecodeString(page)
		if err != nil {
			return nil, "", err
		}

		startPath = string(startPathBytes)
		started = false
	}

	col := collector{
		notes:             make([]string, 0),
		pageSize:          pageSize,
		hasFoundStartPath: started,
	}

	lastPath, err := s.exploreDir(&col, startPath, "/", 0, -1, false)
	if err != nil {
		return nil, "", err
	}

	nextPage := base64.StdEncoding.EncodeToString([]byte(lastPath))

	return col.notes, nextPage, nil
}
