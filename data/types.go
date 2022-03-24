package data

import (
	"encoding/json"
	"time"
)

const DateTimeFormat = "2006-01-02T15:04:05.000000Z07:00"

type DateTime struct {
	time.Time
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt.Format(DateTimeFormat))
}

func (dt *DateTime) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	t, err := time.Parse(DateTimeFormat, str)
	if err != nil {
		return err
	}
	*dt = DateTime{t}
	return nil
}
