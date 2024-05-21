package common

import (
	"database/sql/driver"
	"time"
)

const DbTimeFormat = "2006-01-02T15:04:05"
const JsonTimeFormat = "2006-01-02 15:04:05"

type LocalTime time.Time

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = LocalTime(time.Time{})
		return
	}
	now, err := time.Parse(`"`+DbTimeFormat+`"`, string(data))
	if err != nil {
		now, err = time.Parse(`"`+JsonTimeFormat+`"`, string(data))
	}
	*t = LocalTime(now)
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(JsonTimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, JsonTimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t LocalTime) Value() (driver.Value, error) {
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return []byte(time.Time(t).Format(DbTimeFormat)), nil
}

func (t *LocalTime) Scan(v interface{}) error {
	tTime, _ := time.Parse("2006-01-02T15:04:05Z07:00", v.(time.Time).String())
	*t = LocalTime(tTime)
	return nil
}

func (t LocalTime) String() string {
	return time.Time(t).Format(JsonTimeFormat)
}

func (t LocalTime) DbString() string {
	return time.Time(t).Format(DbTimeFormat)
}
