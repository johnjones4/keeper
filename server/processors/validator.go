package processors

import "main/core"

func (p *Processors) validate(note *core.Note) error {
	if note.Title == "" {
		return &core.ValidationError{
			Message: "title missing",
		}
	}

	if note.Body.Text == "" && len(note.Body.StructuredData) == 0 {
		return &core.ValidationError{
			Message: "body missing",
		}
	}

	return nil
}
