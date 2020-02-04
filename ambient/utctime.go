package ambient

import (
	"time"
)

type utcTime struct {
	time.Time
}

func (u utcTime) format() string {
	return u.Time.UTC().String()
}

func (u utcTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + u.format() + `"`), nil
}
