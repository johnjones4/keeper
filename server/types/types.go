package types

import (
	"github.com/johnjones4/keeper/core"
	"github.com/sirupsen/logrus"
)

type Store interface {
	SaveNote(note *core.Note) error
	RenameAndSaveNote(note *core.Note, oldName string) error
	DeleteNote(note *core.Note) error
	GetNote(id string) (core.Note, error)
	GetNoteByPath(path string) (core.Note, error)
	GetNotes(query core.NotesQuery) ([]core.Note, error)
	GetTags() ([]core.TagInfo, error)
}

type Processor func(note *core.Note) error

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
