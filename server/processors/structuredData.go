package processors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/johnjones4/keeper/core"

	"github.com/PuerkitoBio/goquery"
	"github.com/iand/microdata"
	"github.com/piprate/json-gold/ld"
	"golang.org/x/net/html"
)

func (p *Processors) structuredData(note *core.Note) error {
	if note.Body.Text != "" || note.ID != "" || note.SourceURL == "" {
		return nil
	}

	baseUrl, err := url.Parse(note.SourceURL)
	if err != nil {
		return err
	}

	response, err := http.Get(note.SourceURL)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad http response: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	note.Body.StructuredData = make([]core.StructuredDataProperty, 0)

	err = parseMicrodata(body, baseUrl, note)
	if err != nil {
		return err
	}

	err = parseJsonLd(body, baseUrl, note)
	if err != nil {
		return err
	}

	if len(note.Body.StructuredData) > 0 {
		note.Format = "application/json"
	}

	return nil
}

func parseMicrodata(body []byte, baseUrl *url.URL, note *core.Note) error {
	mdp := microdata.NewParser(bytes.NewBuffer(body), baseUrl)
	data, err := mdp.Parse()
	if err != nil {
		return err
	}

	props, err := normalizeMicrodata(data.Items)
	if err != nil {
		return err
	}
	note.Body.StructuredData = props
	fformatted, _ := json.MarshalIndent(props, "", "  ")
	log.Println(string(fformatted))

	return nil
}

func parseJsonLd(body []byte, baseUrl *url.URL, note *core.Note) error {
	proc := ld.NewJsonLdProcessor()
	opts := ld.NewJsonLdOptions(baseUrl.String())

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
	if err != nil {
		return nil
	}

	sel := doc.Find("[type=\"application/ld+json\"]")
	for _, node := range sel.Nodes {
		s := &goquery.Selection{
			Nodes: []*html.Node{node},
		}

		var mapped []map[string]any
		err := json.Unmarshal([]byte(s.Text()), &mapped)
		if err != nil {
			return err
		}

		for _, d := range mapped {
			expanded, err := proc.Expand(d, opts)
			if err != nil {
				return err
			}
			normalized, err := normalizeJsonLDData(expanded, "")
			if err != nil {
				return err
			}
			note.Body.StructuredData = append(note.Body.StructuredData, normalized...)
		}
	}

	return nil
}

func normalizeJsonLDData(expandedLd []any, keyType string) ([]core.StructuredDataProperty, error) {
	props := make([]core.StructuredDataProperty, 0)
	for _, p := range expandedLd {
		var prop core.StructuredDataProperty
		switch d := p.(type) {
		case map[string]any:
			if _, ok := d["@type"]; ok {
				prop.Type = getStringsValue(d, "@type")
				prop.Properties = make([]core.StructuredDataProperty, 0)
				for k, v := range d {
					if va, ok := v.([]any); ok && k != "@type" {
						props, err := normalizeJsonLDData(va, k)
						if err != nil {
							return nil, err
						}
						prop.Properties = append(prop.Properties, props...)
					}
				}
			} else if value, ok := d["@value"]; ok {
				prop.Type = []string{keyType}
				switch valueT := value.(type) {
				case string:
					prop.String = valueT
				case int:
					prop.Int = valueT
				case float64:
					prop.Float = valueT
				case bool:
					prop.Bool = valueT
				default:
					return nil, fmt.Errorf("no handler for value %s", fmt.Sprint(value))
				}
			} else if ida, ok := d["@id"]; ok {
				prop.Type = []string{keyType}
				switch idT := ida.(type) {
				case string:
					prop.ID = idT
				default:
					return nil, fmt.Errorf("no handler for id %s", fmt.Sprint(ida))
				}
			}
		}
		if prop.Type == nil {
			panic(fmt.Sprint(p))
		}
		props = append(props, prop)
	}
	return props, nil
}

func getStringsValue(m map[string]any, key string) []string {
	v, ok := m[key]
	if !ok {
		return []string{}
	}

	switch vt := v.(type) {
	case []string:
		return vt
	case []any:
		va := make([]string, len(vt))
		for i, ii := range vt {
			va[i] = fmt.Sprint(ii)
		}
		return va
	}

	return []string{}
}

func normalizeMicrodata(items []*microdata.Item) ([]core.StructuredDataProperty, error) {
	props := make([]core.StructuredDataProperty, 0)

	for _, item := range items {
		var prop core.StructuredDataProperty

		prop.Type = item.Types
		prop.ID = item.ID

		prop.Properties = make([]core.StructuredDataProperty, 0)
		for key, val := range item.Properties {
			for _, item := range val {
				switch t := item.(type) {
				case string:
					var p core.StructuredDataProperty
					p.Type = []string{key}
					p.String = t
					prop.Properties = append(prop.Properties, p)
				case *microdata.Item:
					subProps, err := normalizeMicrodata([]*microdata.Item{t})
					if err != nil {
						return nil, err
					}
					prop.Properties = append(prop.Properties, subProps...)
				default:
					return nil, fmt.Errorf("no handler for %s", fmt.Sprint(item))
				}
			}
		}

		props = append(props, prop)
	}

	return props, nil
}
