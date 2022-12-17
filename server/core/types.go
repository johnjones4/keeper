package core

import (
	"time"

	"github.com/sirupsen/logrus"
)

type StructuredDataProperty struct {
	Type       []string                 `json:"type"`
	String     string                   `json:"string,omitempty"`
	Int        int                      `json:"int,omitempty"`
	Float      float64                  `json:"float,omitempty"`
	Bool       bool                     `json:"bool,omitempty"`
	ID         string                   `json:"id,omitempty"`
	Properties []StructuredDataProperty `json:"properties,omitempty"`
}

type NoteBody struct {
	StructuredData []StructuredDataProperty `json:"structuredData,omitempty"`
	Text           string                   `json:"text,omitempty"`
}

type Note struct {
	ID        string    `json:"id"`
	Path      string    `json:"path"`
	Title     string    `json:"title"`
	Body      NoteBody  `json:"body"`
	Tags      []string  `json:"tags"`
	SourceURL string    `json:"sourceURL"`
	Source    string    `json:"source"`
	Format    string    `json:"format"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}

type NotesQuery struct {
	Text  string
	Start time.Time
	End   time.Time
	Tags  []string
}

type TagInfo struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

type Store interface {
	SaveNote(note *Note) error
	RenameAndSaveNote(note *Note, oldName string) error
	DeleteNote(note *Note) error
	GetNote(id string) (Note, error)
	GetNoteByPath(path string) (Note, error)
	GetNotes(query NotesQuery) ([]Note, error)
	GetTags() ([]TagInfo, error)
}

type Processor func(note *Note) error

type Runtime struct {
	Store      Store
	Log        logrus.FieldLogger
	Processors []Processor
}

type ValidationError struct {
	Message string
}

func (v *ValidationError) Error() string {
	return v.Message
}
