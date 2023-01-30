package processors

import (
	"main/types"

	"github.com/johnjones4/keeper/core"
)

func (p *Processors) validate(note *core.Note) error {
	if note.Title == "" {
		return &types.ValidationError{
			Message: "title missing",
		}
	}

	// if note.Body.Text == "" && len(note.Body.StructuredData) == 0 {
	// 	return &core.ValidationError{
	// 		Message: "body missing",
	// 	}
	// }

	return nil
}
