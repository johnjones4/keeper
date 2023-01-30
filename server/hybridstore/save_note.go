package hybridstore

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"time"

	"github.com/johnjones4/keeper/core"

	"github.com/google/uuid"
)

func (s *HybridStore) SaveNote(note *core.Note) error {
	isNew := note.ID == ""
	if isNew {
		note.ID = uuid.NewString()
		note.Created = time.Now()
	}
	note.Updated = time.Now()

	if note.Path != "" {
		fullPath := path.Join(s.filePath, note.Path)

		_, err := os.Stat(fullPath)
		if errors.Is(err, os.ErrNotExist) {
			note.Path = ""
		} else if err != nil {
			return err
		}
	}

	if note.Path == "" {
		filename, err := generateNoteFilename(note.Title, note.Format)
		if err != nil {
			return err
		}

		note.Path = path.Join(prefixZeros(note.Created.Year(), 4), prefixZeros(int(note.Created.Month()), 2), filename)
		fullPath := path.Join(s.filePath, note.Path)

		err = os.MkdirAll(path.Dir(fullPath), 0755)
		if err != nil {
			return err
		}
	}
	fullPath := path.Join(s.filePath, note.Path)

	var contents []byte
	var err error

	if note.Format == "application/json" {
		contents, err = json.MarshalIndent(note.Body.StructuredData, "", "  ")
		if err != nil {
			return err
		}
	} else {
		contents = []byte(note.Body.Text)
	}

	err = os.WriteFile(fullPath, contents, 0755)
	if err != nil {
		return err
	}

	if isNew {
		_, err := s.db.Exec("INSERT INTO notes (id, path, title, sourceURL, source, format, created, updated) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
			note.ID,
			note.Path,
			note.Title,
			note.SourceURL,
			note.Source,
			note.Format,
			note.Created.Unix(),
			note.Updated.Unix(),
		)
		if err != nil {
			return err
		}
	} else {
		_, err := s.db.Exec("UPDATE notes SET path = ?, title = ?, sourceURL = ?, source = ?, format = ?, updated = ? WHERE id = ?",
			note.Path,
			note.Title,
			note.SourceURL,
			note.Source,
			note.Format,
			note.Updated.Unix(),
			note.ID,
		)
		if err != nil {
			return err
		}

		_, err = s.db.Exec(
			"DELETE FROM tags_notes WHERE note_id = $1",
			note.ID,
		)
		if err != nil {
			return err
		}
	}

	note.Tags = uniqueLc(note.Tags)
	for _, tag := range note.Tags {
		_, err := s.db.Exec(
			"INSERT INTO tags_notes (tag, note_id) VALUES ($1, $2)",
			tag,
			note.ID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
