package common

import (
	"database/sql/driver"
	"github.com/guregu/null"
	"time"
)

type LocalNullableTime null.Time

func (t *LocalNullableTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		*t = LocalNullableTime(null.Time{})
		return
	}

	now, err := time.Parse(`"`+DbTimeFormat+`"`, string(data))
	if err != nil {
		now, err = time.Parse(`"`+JsonTimeFormat+`"`, string(data))
	}
	*t = LocalNullableTime(null.TimeFrom(now))
	return
}

func (t LocalNullableTime) MarshalJSON() ([]byte, error) {
	nt := null.Time(t)
	if !nt.Valid {
		return []byte(`null`), nil
	} else {
		return []byte(`"` + nt.Time.Format(JsonTimeFormat) + `"`), nil
	}

}

func (t LocalNullableTime) Value() (driver.Value, error) {
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return []byte(null.Time(t).Time.Format(DbTimeFormat)), nil
}

func (t *LocalNullableTime) Scan(v interface{}) error {
	tTime, _ := time.Parse("2006-01-02T15:04:05Z07:00", v.(time.Time).String())
	*t = LocalNullableTime(null.TimeFrom(tTime))
	return nil
}

func (t LocalNullableTime) String() string {
	return null.Time(t).Time.Format(JsonTimeFormat)
}
func (t LocalNullableTime) DbString() string {
	return null.Time(t).Time.Format(DbTimeFormat)
}
