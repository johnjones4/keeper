package store

import (
	"main/core"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	return c.pageSize != -1 && len(c.notes) >= c.pageSize
}

func (s *StaticFileStore) isIgnored(f string) bool {
	for _, d := range s.ignore {
		if d == f {
			return true
		}
	}
	return false
}

func (s *StaticFileStore) exploreDir(col *collector, startPath string, dir string, depth int, maxDepth int, includeDirs bool) (string, error) {
	if maxDepth != -1 && depth >= maxDepth {
		return "", nil
	}

	realDir := path.Join(s.rootDir, dir)
	if !s.isPathSafe(realDir) {
		return "", core.ErrorBadPath
	}

	contents, err := os.ReadDir(realDir)
	if err != nil {
		return "", err
	}

	for _, file := range contents {
		if !s.isIgnored(file.Name()) {
			keyPath := path.Join(dir, file.Name())
			if includeDirs || !file.IsDir() {
				if col.hasFoundStartPath {
					col.addNote(keyPath)
					if col.isFull() {
						return keyPath, nil
					}
				} else if keyPath == startPath {
					col.hasFoundStartPath = true
				}
			}
			if file.IsDir() {
				lastPath, err := s.exploreDir(col, startPath, keyPath, depth+1, maxDepth, includeDirs)
				if err != nil {
					return "", err
				}
				if lastPath != "" {
					return lastPath, nil
				}
			}
		}
	}

	return "", nil
}

func (s *StaticFileStore) isPathSafe(unsafePath string) bool {
	rel, _ := filepath.Rel(s.rootDir, unsafePath)
	return !strings.Contains(rel, "..")
}
