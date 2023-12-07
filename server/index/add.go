package index

import (
	"main/core"
	"time"
)

func (i *Index) Add(n *core.Note) error {
	_, err := i.db.Exec("INSERT INTO mod_index (keypath, modified) VALUES (?, ?)", n.Key, time.Time(n.Modified).Unix())
	if err != nil {
		return err
	}

	_, err = i.db.Exec("INSERT INTO search_index (keypath, body) VALUES (?, ?)", n.Key, n.Body)
	if err != nil {
		return err
	}

	return nil
}
