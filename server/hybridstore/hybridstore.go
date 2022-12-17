package hybridstore

import (
	"database/sql"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schema string

type HybridStore struct {
	db          *sql.DB
	filePath    string
	recyclePath string
}

func New(dbPath, filePath string, recyclePath string) (*HybridStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	s := HybridStore{
		filePath:    filePath,
		db:          db,
		recyclePath: recyclePath,
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
