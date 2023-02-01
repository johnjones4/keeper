package index

func (i *Index) Search(query string) ([]string, error) {
	rows, err := i.db.Query("SELECT keypath FROM search_index WHERE search_index MATCH ? ORDER BY rank", query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := make([]string, 0)
	for rows.Next() {
		var keypath string
		err = rows.Scan(&keypath)
		if err != nil {
			return nil, err
		}

		items = append(items, keypath)
	}

	return items, nil
}
