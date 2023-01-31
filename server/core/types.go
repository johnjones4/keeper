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

type RuntimeContext struct {
	Store        Store
	PrivateKey   []byte
	PasswordHash []byte
}
