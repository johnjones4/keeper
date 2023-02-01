package core

import "time"

type Note struct {
	Key      string    `json:"key"`
	Body     string    `json:"body"`
	Modified time.Time `json:"modified"`
}

type Store interface {
	SaveNote(n *Note, failOnOverwrite bool, failOnNew bool) error
	GetNote(key string) (Note, error)
	GetNotes(pageSize int, page string) ([]string, string, error)
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
}
