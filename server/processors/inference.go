package processors

import (
	"main/core"
)

func (p *Processors) inference(note *core.Note) error {
	if note.ID != "" || len(note.Body.StructuredData) == 0 {
		return nil
	}

	headline := findProp(note.Body.StructuredData, []string{"http://schema.org/headline", "headline"})
	if headline.String != "" {
		note.Title = headline.String
	}

	for _, prop := range note.Body.StructuredData {
		note.Tags = append(note.Tags, prop.Type...)
	}

	return nil
}

func findProp(props []core.StructuredDataProperty, propNames []string) core.StructuredDataProperty {
	for _, prop := range props {
		for _, propName := range propNames {
			if propsIsOfType(prop, propName) {
				return prop
			}
		}
		if len(prop.Properties) > 0 {
			p := findProp(prop.Properties, propNames)
			if len(p.Type) > 0 {
				return p
			}
		}
	}
	return core.StructuredDataProperty{}
}

func propsIsOfType(prop core.StructuredDataProperty, propName string) bool {
	for _, t := range prop.Type {
		if t == propName {
			return true
		}
	}
	return false
}
