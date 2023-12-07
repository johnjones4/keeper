package core

import (
	"fmt"
	"time"
)

func (t NoteTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339))
	return []byte(stamp), nil
}

func (t *NoteTime) UnmarshalJSON(str []byte) error {
	if len(str) < 3 {
		return nil
	}
	tt, err := time.Parse(time.RFC3339, string(str[1:len(str)-1]))
	if err != nil {
		return err
	}
	*t = NoteTime(tt)
	return nil
}
