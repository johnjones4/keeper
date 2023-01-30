package hybridstore

import "github.com/johnjones4/keeper/core"

func (s *HybridStore) GetTags() ([]core.TagInfo, error) {
	rows, err := s.db.Query("SELECT tag, COUNT(*) FROM tags_notes GROUP BY tag ORDER BY tag")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]core.TagInfo, 0)
	for rows.Next() {
		var tag core.TagInfo
		err = rows.Scan(&tag.Tag, &tag.Count)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
