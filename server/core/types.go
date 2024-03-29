package core

import (
	"time"

	"go.uber.org/zap"
)

type NoteTime time.Time

type Note struct {
	Key      string   `json:"key"`
	Body     string   `json:"body"`
	Modified NoteTime `json:"modified"`
}

type Store interface {
	SaveNote(n *Note, failOnOverwrite bool, failOnNew bool) error
	GetNote(key string) (Note, error)
	GetNotes(pageSize int, page string) ([]string, string, error)
	GetNoteDirectory(directory string) ([]string, error)
}

type Index interface {
	ReIndex() error
	Add(n *Note) error
	Update(n *Note) error
	Search(query string) ([]string, error)
}

type RuntimeContext struct {
	Store        Store
	Index        Index
	PrivateKey   []byte
	PasswordHash []byte
	Log          *zap.SugaredLogger
}
