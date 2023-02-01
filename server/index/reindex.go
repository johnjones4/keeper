package index

import (
	"database/sql"
	"log"
)

func (i *Index) ReIndex() error {
	touchedPaths := make([]string, 0)

	firstPass := true
	nextPage := ""

	for firstPass || nextPage != "" {
		firstPass = false

		log.Printf("Getting notes for page %s", nextPage)

		notes, page, err := i.store.GetNotes(100, nextPage)
		if err != nil {
			return err
		}

		nextPage = page

		for _, keypath := range notes {
			touchedPaths = append(touchedPaths, keypath)

			var modified int
			err := i.db.QueryRow("SELECT modified FROM mod_index WHERE keypath = ?", keypath).Scan(&modified)

			noRows := err != nil && err == sql.ErrNoRows

			if err != nil && !noRows {
				return err
			}

			note, err := i.store.GetNote(keypath)
			if err != nil {
				return err
			}

			if noRows {
				log.Printf("Adding to index: %s", keypath)
				err = i.Add(&note)
				if err != nil {
					return err
				}
			} else if note.Modified.Unix() > int64(modified) {
				log.Printf("Updating index: %s", keypath)
				err = i.Update(&note)
				if err != nil {
					return err
				}
			}
		}
	}

	deletePaths := make([]string, 0)

	rows, err := i.db.Query("SELECT keypath FROM mod_index")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var keypath string
		err = rows.Scan(&keypath)
		if err != nil {
			return err
		}

		if !pathInList(touchedPaths, keypath) {
			log.Printf("Will delete: %s", keypath)
			deletePaths = append(deletePaths, keypath)
		}
	}

	for _, keypath := range deletePaths {
		log.Printf("Deleting: %s", keypath)

		_, err = i.db.Exec("DELETE FROM mod_index WHERE keypath = ?", keypath)
		if err != nil {
			return err
		}

		_, err = i.db.Exec("DELETE FROM search_index WHERE keypath = ?", keypath)
		if err != nil {
			return err
		}
	}

	return nil
}

func pathInList(list []string, path string) bool {
	for _, p := range list {
		if p == path {
			return true
		}
	}
	return false
}
