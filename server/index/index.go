package index

import (
	"database/sql"
	"main/core"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

//go:embed schema.sql
var schema string

type Index struct {
	store core.Store
	log   *zap.SugaredLogger
	db    *sql.DB
}

func New(filepath string, store core.Store, log *zap.SugaredLogger) (*Index, error) {
	s := Index{
		store: store,
		log:   log,
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
