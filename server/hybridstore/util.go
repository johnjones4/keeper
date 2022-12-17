package hybridstore

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"main/core"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/flytam/filenamify"
)

type scannable interface {
	Scan(args ...any) error
}

func parseRow(s scannable) (core.Note, error) {
	var createdUnix int64
	var updatedUnix int64
	var note core.Note

	err := s.Scan(
		&note.ID,
		&note.Path,
		&note.Title,
		&note.SourceURL,
		&note.Source,
		&note.Format,
		&createdUnix,
		&updatedUnix,
	)
	if err != nil {
		return core.Note{}, err
	}

	note.Created = time.Unix(createdUnix, 0)
	note.Updated = time.Unix(updatedUnix, 0)

	return note, nil
}

func generateNoteFilename(name, format string) (string, error) {
	safename, err := filenamify.Filenamify(name, filenamify.Options{})
	if err != nil {
		return "", err
	}
	ext := "txt"
	slashIndex := strings.Index(format, "/")
	if slashIndex > 0 {
		ext = format[slashIndex+1:]
	}
	return fmt.Sprintf("%s.%s", safename, ext), nil
}

func (s *HybridStore) loadNoteBody(note *core.Note) error {
	fullPath := path.Join(s.filePath, note.Path)
	fileBytes, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	if note.Format == "application/json" {
		err = json.Unmarshal(fileBytes, &note.Body.StructuredData)
		if err != nil {
			return err
		}
	} else {
		note.Body.Text = string(fileBytes)
	}

	return nil
}

func prefixZeros(n int, length int) string {
	ns := strconv.Itoa(n)
	if len(ns) == length {
		return ns
	} else if len(ns) > length {
		panic(fmt.Sprintf("cannot prefix \"%s\" to length %d", ns, length))
	}
	for i := 0; i < length-len(ns); i++ {
		ns = "0" + ns
	}
	return ns
}

func (s *HybridStore) populateTags(asset *core.Note) error {
	rows, err := s.db.Query("SELECT tag FROM tags_notes WHERE note_id = $1", asset.ID)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}
	defer rows.Close()
	asset.Tags = make([]string, 0)
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return err
		}
		asset.Tags = append(asset.Tags, tag)
	}
	return nil
}

func uniqueLc(arr []string) []string {
	occurred := map[string]bool{}
	result := []string{}
	for _, i := range arr {
		ilc := strings.ToLower(i)
		if !occurred[ilc] {
			occurred[ilc] = true
			result = append(result, ilc)
		}
	}
	return result
}
