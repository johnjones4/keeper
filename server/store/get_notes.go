package store

import (
	"encoding/base64"
	"os"
	"path"
)

type collector struct {
	notes             []string
	pageSize          int
	hasFoundStartPath bool
}

func (c *collector) addNote(path string) {
	c.notes = append(c.notes, path)
}

func (c *collector) isFull() bool {
	return len(c.notes) >= c.pageSize
}

func (s *StaticFileStore) GetNotes(pageSize int, page string) ([]string, string, error) {
	startPath := ""
	started := true
	if page != "" {
		startPathBytes, err := base64.StdEncoding.DecodeString(page)
		if err != nil {
			return nil, "", err
		}

		startPath = string(startPathBytes)
		started = false
	}

	col := collector{
		notes:             make([]string, 0),
		pageSize:          pageSize,
		hasFoundStartPath: started,
	}

	lastPath, err := s.exploreDir(&col, startPath, "/")
	if err != nil {
		return nil, "", err
	}

	nextPage := base64.StdEncoding.EncodeToString([]byte(lastPath))

	return col.notes, nextPage, nil
}

func (s *StaticFileStore) isIgnored(f string) bool {
	for _, d := range s.ignore {
		if d == f {
			return true
		}
	}
	return false
}

func (s *StaticFileStore) exploreDir(col *collector, startPath string, dir string) (string, error) {
	contents, err := os.ReadDir(path.Join(s.rootDir, dir))
	if err != nil {
		return "", err
	}

	for _, file := range contents {
		if !s.isIgnored(file.Name()) {
			keyPath := path.Join(dir, file.Name())
			if file.IsDir() {
				lastPath, err := s.exploreDir(col, startPath, keyPath)
				if err != nil {
					return "", err
				}
				if lastPath != "" {
					return lastPath, nil
				}
			} else {
				if col.hasFoundStartPath {
					col.addNote(keyPath)
					if col.isFull() {
						return keyPath, nil
					}
				} else if keyPath == startPath {
					col.hasFoundStartPath = true
				}
			}
		}
	}

	return "", nil
}
