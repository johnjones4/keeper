package index

import (
	"database/sql"
	"main/core"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schema string

type Index struct {
	store core.Store
	db    *sql.DB
}

func New(filepath string, store core.Store) (*Index, error) {
	s := Index{
		store: store,
	}

	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	s.db = db

	_, err = s.db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
