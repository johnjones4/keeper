package hybridstore

import (
	"main/core"
	"strings"

	"github.com/jmoiron/sqlx"
)

func (s *HybridStore) GetNotes(query core.NotesQuery) ([]core.Note, error) {
	predicates := make([]string, 0)
	params := make([]any, 0)

	if query.Text != "" {
		predicates = append(predicates, "title LIKE '%' || ? || '%'")
		params = append(params, query.Text)
	}

	if !query.Start.IsZero() {
		predicates = append(predicates, "created <= ?")
		params = append(params, query.Start.Unix())
	}

	if !query.End.IsZero() {
		predicates = append(predicates, "created <= ?")
		params = append(params, query.End.Unix())
	}

	if len(query.Tags) > 0 {
		pred, args, err := sqlx.In("id IN (SELECT note_id FROM tags_notes WHERE tag IN (?))", query.Tags)
		if err != nil {
			return nil, err
		}
		predicates = append(predicates, pred)
		params = append(params, args...)
	}

	queryStr := "SELECT id, path, title, sourceURL, source, format, created, updated FROM NOTES"
	if len(predicates) > 0 {
		queryStr += " WHERE " + strings.Join(predicates, " AND ")
	}
	queryStr += " ORDER BY created"

	rows, err := s.db.Query(queryStr, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]core.Note, 0)
	for rows.Next() {
		note, err := parseRow(rows)
		if err != nil {
			return nil, err
		}

		err = s.populateTags(&note)
		if err != nil {
			return nil, err
		}

		err = s.loadNoteBody(&note)
		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	return notes, nil
}
