package index

import (
	"main/core"
	"time"
)

func (i *Index) Update(n *core.Note) error {
	_, err := i.db.Exec("UPDATE mod_index SET modified = ? WHERE keypath = ?", time.Time(n.Modified).Unix(), n.Key)
	if err != nil {
		return err
	}

	_, err = i.db.Exec("UPDATE search_index SET body = ? WHERE keypath = ?", n.Body, n.Key)
	if err != nil {
		return err
	}

	return nil
}
